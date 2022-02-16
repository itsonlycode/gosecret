package action

import (
	"bytes"
	"context"
	"os"
	"testing"

	"github.com/itsonlycode/gosecret/internal/out"
	"github.com/itsonlycode/gosecret/pkg/ctxutil"
	"github.com/itsonlycode/gosecret/tests/gptest"

	_ "github.com/itsonlycode/gosecret/internal/backend/crypto"
	_ "github.com/itsonlycode/gosecret/internal/backend/storage"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v2"
)

func TestVersion(t *testing.T) {
	u := gptest.NewUnitTester(t)
	defer u.Remove()

	ctx := context.Background()
	ctx = ctxutil.WithAlwaysYes(ctx, true)
	ctx = ctxutil.WithInteractive(ctx, false)

	act, err := newMock(ctx, u)
	require.NoError(t, err)

	buf := &bytes.Buffer{}
	out.Stdout = buf
	stdout = buf
	defer func() {
		out.Stdout = os.Stdout
		stdout = os.Stdout
	}()

	cli.VersionPrinter = func(*cli.Context) {
		out.Printf(ctx, "gosecret version 0.0.0-test")
	}

	t.Run("print fixed version", func(t *testing.T) {
		assert.NoError(t, act.Version(gptest.CliCtx(ctx, t)))
	})
}
