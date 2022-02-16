package root

import (
	"context"
	"testing"

	"github.com/itsonlycode/gosecret/pkg/ctxutil"
	"github.com/itsonlycode/gosecret/tests/gptest"

	"github.com/fatih/color"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRecipients(t *testing.T) {
	u := gptest.NewUnitTester(t)
	defer u.Remove()

	ctx := context.Background()
	ctx = ctxutil.WithAlwaysYes(ctx, true)
	ctx = ctxutil.WithHidden(ctx, true)
	color.NoColor = true

	rs, err := createRootStore(ctx, u)
	require.NoError(t, err)

	assert.Equal(t, []string{"0xDEADBEEF"}, rs.ListRecipients(ctx, ""))
	rt, err := rs.RecipientsTree(ctx, false)
	require.NoError(t, err)
	assert.Equal(t, "gosecret\n└── 0xDEADBEEF\n", rt.Format(0))
}
