package cli

import (
	"bytes"
	"context"
	"os"
	"os/exec"

	"github.com/itsonlycode/gosecret/internal/backend/crypto/gpg"
	"github.com/itsonlycode/gosecret/internal/out"
	"github.com/itsonlycode/gosecret/pkg/debug"
)

// Encrypt will encrypt the given content for the recipients. If alwaysTrust is true
// the trust-model will be set to always as to avoid (annoying) "unusable public key"
// errors when encrypting.
func (g *GPG) Encrypt(ctx context.Context, plaintext []byte, recipients []string) ([]byte, error) {
	args := append(g.args, "--encrypt")
	if gpg.IsAlwaysTrust(ctx) {
		// changing the trustmodel is possibly dangerous. A user should always
		// explicitly opt-in to do this
		args = append(args, "--trust-model=always")
	}
	for _, r := range recipients {
		kl, err := g.listKeys(ctx, "public", r)
		if err != nil {
			debug.Log("Failed to check key %s. Adding anyway. %s", err)
		} else if len(kl.UseableKeys(gpg.IsAlwaysTrust(ctx))) < 1 {
			out.Printf(ctx, "Not using expired key %s for encryption", r)
			continue
		}
		args = append(args, "--recipient", r)
	}

	buf := &bytes.Buffer{}

	cmd := exec.CommandContext(ctx, g.binary, args...)
	cmd.Stdin = bytes.NewReader(plaintext)
	cmd.Stdout = buf
	cmd.Stderr = os.Stderr

	debug.Log("%s %+v", cmd.Path, cmd.Args)
	err := cmd.Run()
	return buf.Bytes(), err
}
