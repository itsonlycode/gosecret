package action

import (
	"context"
	"flag"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/itsonlycode/gosecret/internal/backend"
	"github.com/itsonlycode/gosecret/internal/backend/crypto/plain"
	"github.com/itsonlycode/gosecret/internal/config"
	"github.com/itsonlycode/gosecret/pkg/ctxutil"
	"github.com/itsonlycode/gosecret/tests/gptest"

	"github.com/blang/semver/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v2"
)

func newMock(ctx context.Context, u *gptest.Unit) (*Action, error) {
	cfg := config.Load()
	cfg.Path = u.StoreDir("")

	ctx = backend.WithCryptoBackend(ctx, backend.Plain)
	ctx = backend.WithStorageBackend(ctx, backend.GitFS)
	act, err := newAction(cfg, semver.Version{}, false)
	if err != nil {
		return nil, err
	}

	fs := flag.NewFlagSet("default", flag.ContinueOnError)
	c := cli.NewContext(cli.NewApp(), fs, nil)
	c.Context = ctx
	if err := act.IsInitialized(c); err != nil {
		return nil, err
	}

	return act, nil
}

func TestAction(t *testing.T) {
	u := gptest.NewUnitTester(t)
	defer u.Remove()

	ctx := context.Background()
	ctx = ctxutil.WithInteractive(ctx, false)

	act, err := newMock(ctx, u)
	require.NoError(t, err)
	require.NotNil(t, act)

	actName := "action.test"

	if runtime.GOOS == "windows" {
		actName = "action.test.exe"
	}

	assert.Equal(t, actName, act.Name)

	assert.Contains(t, act.String(), u.StoreDir(""))
	assert.Equal(t, 0, len(act.Store.Mounts()))
}

func TestNew(t *testing.T) {
	td, err := os.MkdirTemp("", "gosecret-")
	require.NoError(t, err)
	defer func() {
		_ = os.RemoveAll(td)
	}()

	cfg := config.New()
	sv := semver.Version{}

	t.Run("init a new store", func(t *testing.T) {
		_, err = New(cfg, sv)
		require.NoError(t, err)
	})

	t.Run("init an existing plain store", func(t *testing.T) {
		cfg.Path = filepath.Join(td, "store")
		assert.NoError(t, os.MkdirAll(cfg.Path, 0700))
		assert.NoError(t, os.WriteFile(filepath.Join(cfg.Path, plain.IDFile), []byte("foobar"), 0600))
		_, err = New(cfg, sv)
		assert.NoError(t, err)
	})
}
