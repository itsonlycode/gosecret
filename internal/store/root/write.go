package root

import (
	"context"

	"github.com/itsonlycode/gosecret/pkg/gosecret"
)

// Set encodes and write the ciphertext of one entry to disk
func (r *Store) Set(ctx context.Context, name string, sec gosecret.Byter) error {
	store, name := r.getStore(name)
	return store.Set(ctx, name, sec)
}
