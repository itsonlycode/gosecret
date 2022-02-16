package action

import (
	"bytes"
	"context"
	"os"
	"testing"

	"github.com/itsonlycode/gosecret/internal/out"
	"github.com/itsonlycode/gosecret/pkg/ctxutil"
	"github.com/itsonlycode/gosecret/pkg/gosecret/secrets"
	"github.com/itsonlycode/gosecret/tests/gptest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAudit(t *testing.T) {
	u := gptest.NewUnitTester(t)
	defer u.Remove()

	ctx := context.Background()
	ctx = ctxutil.WithAlwaysYes(ctx, true)
	ctx = ctxutil.WithHidden(ctx, true)
	act, err := newMock(ctx, u)
	require.NoError(t, err)
	require.NotNil(t, act)

	buf := &bytes.Buffer{}
	out.Stdout = buf
	stdout = buf
	defer func() {
		out.Stdout = os.Stdout
		stdout = os.Stdout
	}()

	t.Run("expect audit complaints on very weak passwords", func(t *testing.T) {
		sec := &secrets.Plain{}
		sec.SetPassword("123")
		assert.NoError(t, act.Store.Set(ctx, "bar", sec))
		assert.NoError(t, act.Store.Set(ctx, "baz", sec))

		assert.Error(t, act.Audit(gptest.CliCtx(ctx, t)))
		buf.Reset()
	})

	t.Run("test with filter and very passwords", func(t *testing.T) {
		c := gptest.CliCtx(ctx, t, "foo")
		assert.Error(t, act.Audit(c))
		buf.Reset()
	})

	t.Run("test empty store", func(t *testing.T) {
		for _, v := range []string{"foo", "bar", "baz"} {
			assert.NoError(t, act.Store.Delete(ctx, v))
		}
		assert.NoError(t, act.Audit(gptest.CliCtx(ctx, t)))
		assert.Contains(t, "No secrets found", buf.String())
		buf.Reset()
	})
}
