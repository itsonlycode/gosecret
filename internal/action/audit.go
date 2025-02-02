package action

import (
	"github.com/itsonlycode/gosecret/internal/audit"
	"github.com/itsonlycode/gosecret/internal/out"
	"github.com/itsonlycode/gosecret/internal/tree"
	"github.com/itsonlycode/gosecret/pkg/ctxutil"
	"github.com/itsonlycode/gosecret/pkg/debug"

	"github.com/urfave/cli/v2"
)

// Audit validates passwords against common flaws
func (s *Action) Audit(c *cli.Context) error {
	s.rem.Reset("audit")

	filter := c.Args().First()
	ctx := ctxutil.WithGlobalFlags(c)

	out.Print(ctx, "Auditing passwords for common flaws ...")
	t, err := s.Store.Tree(ctx)
	if err != nil {
		return ExitError(ExitList, err, "failed to get store tree: %s", err)
	}

	if filter != "" {
		subtree, err := t.FindFolder(filter)
		if err != nil {
			return ExitError(ExitUnknown, err, "failed to find subtree: %s", err)
		}
		debug.Log("subtree for %q: %+v", filter, subtree)
		t = subtree
	}
	list := t.List(tree.INF)

	if len(list) < 1 {
		out.Printf(ctx, "No secrets found")
		return nil
	}

	return audit.Batch(ctx, list, s.Store)
}
