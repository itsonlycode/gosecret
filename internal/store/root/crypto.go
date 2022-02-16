package root

import (
	"context"

	"github.com/itsonlycode/gosecret/internal/backend"
	"github.com/itsonlycode/gosecret/pkg/debug"
)

// Crypto returns the crypto backend
func (r *Store) Crypto(ctx context.Context, name string) backend.Crypto {
	sub, _ := r.getStore(name)
	if !sub.Valid() {
		debug.Log("Sub-Store not found for %s. Returning nil crypto backend", name)
		return nil
	}
	return sub.Crypto()
}
