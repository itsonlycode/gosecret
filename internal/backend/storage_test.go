package backend

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/itsonlycode/gosecret/pkg/ctxutil"
	"github.com/itsonlycode/gosecret/pkg/debug"
	"github.com/itsonlycode/gosecret/tests/gptest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDetectStorage(t *testing.T) {
	ctx := context.Background()

	uv := gptest.UnsetVars("GOPASS_HOMEDIR")
	defer uv()

	td, err := os.MkdirTemp("", "gosecret-")
	require.NoError(t, err)
	defer func() {
		_ = os.RemoveAll(td)
	}()

	// all tests involving age should set GOPASS_HOMEDIR
	os.Setenv("GOPASS_HOMEDIR", td)
	ctx = ctxutil.WithPasswordCallback(ctx, func(_ string, _ bool) ([]byte, error) {
		debug.Log("static test password callback")
		return []byte("gosecret"), nil
	})

	fsDir := filepath.Join(td, "fs")
	assert.NoError(t, os.MkdirAll(fsDir, 0700))

	t.Run("detect fs", func(t *testing.T) {
		r, err := DetectStorage(ctx, fsDir)
		assert.NoError(t, err)
		assert.NotNil(t, r)
		assert.Equal(t, "fs", r.Name())
	})
}
