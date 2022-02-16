package root

import (
	"context"

	"github.com/itsonlycode/gosecret/pkg/gosecret"
)

// Get returns the plaintext of a single key
func (r *Store) Get(ctx context.Context, name string) (gosecret.Secret, error) {
	// forward to substore
	store, name := r.getStore(name)
	sec, err := store.Get(ctx, name)
	return sec, err
}
