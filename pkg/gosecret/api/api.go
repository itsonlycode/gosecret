package api

import (
	"context"
	"fmt"

	"github.com/itsonlycode/gosecret/internal/tree"

	_ "github.com/itsonlycode/gosecret/internal/backend/crypto"  // load crypto backends
	_ "github.com/itsonlycode/gosecret/internal/backend/storage" // load storage backends
	"github.com/itsonlycode/gosecret/internal/config"
	"github.com/itsonlycode/gosecret/internal/queue"
	"github.com/itsonlycode/gosecret/internal/store/root"
	"github.com/itsonlycode/gosecret/pkg/gosecret"
)

// Gosecret is a secret store implementation
type Gosecret struct {
	rs *root.Store
}

// make sure that *Gosecret implements Store
var _ gosecret.Store = &Gosecret{}

// New creates a new secret store
// WARNING: This will need to change to accommodate for runtime configuration.
func New(ctx context.Context) (*Gosecret, error) {
	cfg := config.LoadWithFallbackRelaxed()
	store := root.New(cfg)
	initialized, err := store.IsInitialized(ctx)
	if err != nil {
		return nil, err
	}
	if !initialized {
		return nil, fmt.Errorf("store not initialized. run gosecret init first")
	}
	return &Gosecret{
		rs: store,
	}, nil
}

// List returns a list of all secrets.
func (g *Gosecret) List(ctx context.Context) ([]string, error) {
	return g.rs.List(ctx, tree.INF)
}

// Get returns a single, encrypted secret. It must be unwrapped before use.
func (g *Gosecret) Get(ctx context.Context, name, revision string) (gosecret.Secret, error) {
	return g.rs.Get(ctx, name)
}

// Set adds a new revision to an existing secret or creates a new one.
func (g *Gosecret) Set(ctx context.Context, name string, sec gosecret.Byter) error {
	return g.rs.Set(ctx, name, sec)
}

// Remove removes a single secret.
func (g *Gosecret) Remove(ctx context.Context, name string) error {
	return g.rs.Delete(ctx, name)
}

// RemoveAll removes all secrets with a given prefix.
func (g *Gosecret) RemoveAll(ctx context.Context, prefix string) error {
	return g.rs.Prune(ctx, prefix)
}

// Rename move a prefix to another.
func (g *Gosecret) Rename(ctx context.Context, src, dest string) error {
	return g.rs.Move(ctx, src, dest)
}

// Sync synchronizes a secret with a remote
func (g *Gosecret) Sync(ctx context.Context) error {
	return fmt.Errorf("not yet implemented")
}

// Revisions lists all revisions of this secret
func (g *Gosecret) Revisions(ctx context.Context, name string) ([]string, error) {
	return nil, fmt.Errorf("not yet implemented")
}

func (g *Gosecret) String() string {
	return "gosecret"
}

// Close shuts down all background processes
func (g *Gosecret) Close(ctx context.Context) error {
	return queue.GetQueue(ctx).Wait(ctx)
}

// ConfigDir returns gosecret' configuration directory
func ConfigDir() string {
	return config.Directory()
}
