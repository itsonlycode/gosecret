package root

import (
	"context"
	"testing"

	"github.com/fatih/color"
	"github.com/itsonlycode/gosecret/pkg/ctxutil"
	"github.com/itsonlycode/gosecret/tests/gptest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRCS(t *testing.T) {
	u := gptest.NewUnitTester(t)
	defer u.Remove()

	ctx := context.Background()
	ctx = ctxutil.WithAlwaysYes(ctx, true)
	ctx = ctxutil.WithHidden(ctx, true)
	color.NoColor = true

	rs, err := createRootStore(ctx, u)
	require.NoError(t, err)

	assert.NoError(t, rs.RCSStatus(ctx, ""))

	revs, err := rs.ListRevisions(ctx, "foo")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(revs))
}
