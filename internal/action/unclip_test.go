package action

import (
	"bytes"
	"context"
	"os"
	"testing"

	"github.com/itsonlycode/gosecret/internal/out"
	"github.com/itsonlycode/gosecret/tests/gptest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	_ "github.com/itsonlycode/gosecret/internal/backend/crypto"
	_ "github.com/itsonlycode/gosecret/internal/backend/storage"
)

func TestUnclip(t *testing.T) {
	u := gptest.NewUnitTester(t)
	defer u.Remove()

	buf := &bytes.Buffer{}
	out.Stdout = buf
	stdout = buf
	defer func() {
		out.Stdout = os.Stdout
		stdout = os.Stdout
	}()

	ctx := context.Background()
	act, err := newMock(ctx, u)
	require.NoError(t, err)
	require.NotNil(t, act)

	t.Run("unlcip should fail", func(t *testing.T) {
		assert.Error(t, act.Unclip(gptest.CliCtxWithFlags(ctx, t, map[string]string{"timeout": "0"})))
	})
}
