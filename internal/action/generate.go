package action

import (
	"context"
	"fmt"
	"path"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/itsonlycode/gosecret/internal/tree"
	"github.com/itsonlycode/gosecret/pkg/gosecret/secrets"

	"github.com/itsonlycode/gosecret/internal/out"
	"github.com/itsonlycode/gosecret/pkg/clipboard"
	"github.com/itsonlycode/gosecret/pkg/ctxutil"
	"github.com/itsonlycode/gosecret/pkg/debug"
	"github.com/itsonlycode/gosecret/pkg/gosecret"
	"github.com/itsonlycode/gosecret/pkg/pwgen/pwrules"
	"github.com/itsonlycode/gosecret/pkg/termio"

	"github.com/urfave/cli/v2"
)

var (
	reNumber = regexp.MustCompile(`^\d+$`)
)

// Generate and save a password
func (s *Action) Generate(c *cli.Context) error {
	ctx := ctxutil.WithGlobalFlags(c)
	ctx = WithClip(ctx, c.Bool("clip"))
	force := c.Bool("force")
	edit := c.Bool("edit")

	args, kvps := parseArgs(c)
	name := args.Get(0)
	key, length := keyAndLength(args)

	ctx = ctxutil.WithForce(ctx, force)

	// ask for name of the secret if it wasn't provided already
	if name == "" {
		var err error
		name, err = termio.AskForString(ctx, "Which name do you want to use?", "")
		if err != nil || name == "" {
			return ExitError(ExitNoName, err, "please provide a password name")
		}
	}

	// ask for confirmation before overwriting existing entry
	if !force { // don't check if it's force anyway
		if s.Store.Exists(ctx, name) && key == "" && !termio.AskForConfirmation(ctx, fmt.Sprintf("An entry already exists for %s. Overwrite the current password?", name)) {
			return ExitError(ExitAborted, nil, "user aborted. not overwriting your current password")
		}
	}

	// generate password
	password, err := s.generatePassword(ctx, c, length, name)
	if err != nil {
		return err
	}

	// display or copy to clipboard
	if err := s.generateCopyOrPrint(ctx, c, name, key, password); err != nil {
		return err
	}

	// write generated password to store
	ctx, err = s.generateSetPassword(ctx, name, key, password, kvps)
	if err != nil {
		return err
	}

	// if requested launch editor to add more data to the generated secret
	if edit && termio.AskForConfirmation(ctx, fmt.Sprintf("Do you want to add more data for %s?", name)) {
		c.Context = ctx
		if err := s.Edit(c); err != nil {
			return ExitError(ExitUnknown, err, "failed to edit %q: %s", name, err)
		}
	}

	return nil
}

func keyAndLength(args argList) (string, string) {
	key := args.Get(1)
	length := args.Get(2)

	// generate can be called with one positional arg or two
	// one - the desired length for the "master" secret itself
	// two - the key in a YAML doc and the length for a secret generated for this
	// key only
	if length == "" && key != "" && reNumber.MatchString(key) {
		length = key
		key = ""
	}

	return key, length
}

// generateCopyOrPrint will print the password to the screen or copy to the
// clipboard
func (s *Action) generateCopyOrPrint(ctx context.Context, c *cli.Context, name, key, password string) error {
	entry := name
	if key != "" {
		entry += ":" + key
	}

	out.OKf(ctx, "Password for entry %q generated", entry)

	// copy to clipboard if:
	// - explicitly requested with -c
	// - autoclip=true, but only if output is not being redirected
	if IsClip(ctx) || (s.cfg.AutoClip && ctxutil.IsTerminal(ctx)) {
		if err := clipboard.CopyTo(ctx, name, []byte(password), s.cfg.ClipTimeout); err != nil {
			return ExitError(ExitIO, err, "failed to copy to clipboard: %s", err)
		}
		// if autoclip is on and we're not printing the password to the terminal
		// at least leave a notice that we did indeed copy it
		if s.cfg.AutoClip && !c.Bool("print") {
			out.Print(ctx, "Copied to clipboard")
			return nil
		}
	}

	if !c.Bool("print") {
		out.Printf(ctx, "Not printing secrets by default. Use 'gosecret show %s' to display the password.", entry)
		return nil
	}
	if c.IsSet("print") && !c.Bool("print") && ctxutil.IsShowSafeContent(ctx) {
		debug.Log("safecontent suppresing printing")
		return nil
	}

	out.Printf(
		ctx,
		"⚠ The generated password is:\n\n%s\n",
		out.Secret(password),
	)
	return nil
}

func hasPwRuleForSecret(name string) (string, pwrules.Rule) {
	for name != "" && name != "." {
		d := path.Base(name)
		if r, found := pwrules.LookupRule(d); found {
			return d, r
		}
		name = path.Dir(name)
	}
	return "", pwrules.Rule{}
}

// generateSetPassword will update or create a secret
func (s *Action) generateSetPassword(ctx context.Context, name, key, password string, kvps map[string]string) (context.Context, error) {
	// set a single key in an entry
	if key != "" {
		sec, err := s.Store.Get(ctx, name)
		if err != nil {
			return ctx, ExitError(ExitEncrypt, err, "failed to set key %q of %q: %s", key, name, err)
		}
		setMetadata(sec, kvps)
		sec.Set(key, password)
		if err := s.Store.Set(ctxutil.WithCommitMessage(ctx, "Generated password for key"), name, sec); err != nil {
			return ctx, ExitError(ExitEncrypt, err, "failed to set key %q of %q: %s", key, name, err)
		}
		return ctx, nil
	}

	// replace password in existing secret
	if s.Store.Exists(ctx, name) {
		ctx, err := s.generateReplaceExisting(ctx, name, key, password, kvps)
		if err == nil {
			return ctx, nil
		}
		out.Errorf(ctx, "Failed to read existing secret. Creating anew. Error: %s", err.Error())
	}

	// generate a completely new secret
	var sec gosecret.Secret
	sec = secrets.New()
	sec.SetPassword(password)
	if u := hasChangeURL(name); u != "" {
		sec.Set("password-change-url", u)
	}

	if content, found := s.renderTemplate(ctx, name, []byte(password)); found {
		nSec := &secrets.Plain{}
		if _, err := nSec.Write(content); err == nil {
			sec = nSec
		} else {
			debug.Log("failed to handle template: %s", err)
		}
	}

	if err := s.Store.Set(ctxutil.WithCommitMessage(ctx, "Generated Password"), name, sec); err != nil {
		return ctx, ExitError(ExitEncrypt, err, "failed to create %q: %s", name, err)
	}
	return ctx, nil
}

func hasChangeURL(name string) string {
	p := strings.Split(name, "/")
	for i := len(p) - 1; i > 0; i-- {
		if u := pwrules.LookupChangeURL(p[i]); u != "" {
			return u
		}
	}
	return ""
}

func (s *Action) generateReplaceExisting(ctx context.Context, name, key, password string, kvps map[string]string) (context.Context, error) {
	sec, err := s.Store.Get(ctx, name)
	if err != nil {
		return ctx, ExitError(ExitEncrypt, err, "failed to set key %q of %q: %s", key, name, err)
	}

	setMetadata(sec, kvps)
	sec.SetPassword(password)
	if err := s.Store.Set(ctxutil.WithCommitMessage(ctx, "Generated password for YAML key"), name, sec); err != nil {
		return ctx, ExitError(ExitEncrypt, err, "failed to set key %q of %q: %s", key, name, err)
	}

	return ctx, nil
}

func setMetadata(sec gosecret.Secret, kvps map[string]string) {
	for k, v := range kvps {
		sec.Set(k, v)
	}
}

// CompleteGenerate implements the completion heuristic for the generate command
func (s *Action) CompleteGenerate(c *cli.Context) {
	ctx := ctxutil.WithGlobalFlags(c)
	if c.Args().Len() < 1 {
		return
	}
	needle := c.Args().Get(0)

	_, err := s.Store.IsInitialized(ctx) // important to make sure the structs are not nil
	if err != nil {
		out.Errorf(ctx, "Store not initialized: %s", err)
		return
	}
	list, err := s.Store.List(ctx, tree.INF)
	if err != nil {
		return
	}

	if strings.Contains(needle, "/") {
		list = filterPrefix(uniq(extractEmails(list)), path.Base(needle))
	} else {
		list = filterPrefix(uniq(extractDomains(list)), needle)
	}

	for _, v := range list {
		fmt.Fprintln(stdout, bashEscape(v))
	}
}

func extractEmails(list []string) []string {
	results := make([]string, 0, len(list))
	for _, e := range list {
		e = path.Base(e)
		if strings.Contains(e, "@") || strings.Contains(e, "_") {
			results = append(results, e)
		}
	}
	return results
}

var reDomain = regexp.MustCompile(`^(?i)([a-z0-9]+(-[a-z0-9]+)*\.)+[a-z]{2,}$`)

func extractDomains(list []string) []string {
	results := make([]string, 0, len(list))
	for _, e := range list {
		e = path.Base(e)
		if reDomain.MatchString(e) {
			results = append(results, e)
		}
	}
	return results
}

func uniq(in []string) []string {
	set := make(map[string]struct{}, len(in))
	for _, e := range in {
		set[e] = struct{}{}
	}
	out := make([]string, 0, len(set))
	for k := range set {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

func filterPrefix(in []string, prefix string) []string {
	out := make([]string, 0, len(in))
	for _, e := range in {
		if strings.HasPrefix(e, prefix) {
			out = append(out, e)
		}
	}
	return out
}
