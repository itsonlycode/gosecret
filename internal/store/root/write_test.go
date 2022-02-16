package root

import (
	"context"
	"testing"

	"github.com/itsonlycode/gosecret/pkg/ctxutil"
	"github.com/itsonlycode/gosecret/pkg/gosecret/secrets"
	"github.com/itsonlycode/gosecret/tests/gptest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSet(t *testing.T) {
	u := gptest.NewUnitTester(t)
	defer u.Remove()

	ctx := context.Background()
	ctx = ctxutil.WithAlwaysYes(ctx, true)
	ctx = ctxutil.WithHidden(ctx, true)

	rs, err := createRootStore(ctx, u)
	require.NoError(t, err)

	sec := &secrets.Plain{}
	sec.SetPassword("foo")
	sec.WriteString("bar")
	assert.NoError(t, rs.Set(ctx, "zab", sec))

	err = rs.Set(ctx, "zab2", sec)
	assert.NoError(t, err)
}
