package leaf

import (
	"context"
	"fmt"

	"github.com/itsonlycode/gosecret/internal/backend"
	"github.com/itsonlycode/gosecret/pkg/ctxutil"
)

func (s *Store) initStorageBackend(ctx context.Context) error {
	ctx = ctxutil.WithAlias(ctx, s.alias)
	store, err := backend.DetectStorage(ctx, s.path)
	if err != nil {
		return fmt.Errorf("unknown storage backend: %w", err)
	}
	s.storage = store
	return nil
}
