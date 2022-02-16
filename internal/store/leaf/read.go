package leaf

import (
	"context"

	"github.com/itsonlycode/gosecret/pkg/ctxutil"
	"github.com/itsonlycode/gosecret/pkg/gosecret/secrets"

	"github.com/itsonlycode/gosecret/internal/out"
	"github.com/itsonlycode/gosecret/internal/store"
	"github.com/itsonlycode/gosecret/pkg/debug"
	"github.com/itsonlycode/gosecret/pkg/gosecret"
	"github.com/itsonlycode/gosecret/pkg/gosecret/secrets/secparse"
)

// Get returns the plaintext of a single key
func (s *Store) Get(ctx context.Context, name string) (gosecret.Secret, error) {
	p := s.passfile(name)

	ciphertext, err := s.storage.Get(ctx, p)
	if err != nil {
		debug.Log("File %s not found: %s", p, err)
		return nil, store.ErrNotFound
	}

	content, err := s.crypto.Decrypt(ctx, ciphertext)
	if err != nil {
		out.Errorf(ctx, "Decryption failed: %s\n%s", err, string(content))
		return nil, store.ErrDecrypt
	}

	if !ctxutil.IsShowParsing(ctx) {
		return secrets.ParsePlain(content), nil
	}

	return secparse.Parse(content)
}
