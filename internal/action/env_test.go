package action

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"runtime"
	"testing"

	"github.com/itsonlycode/gosecret/internal/out"
	"github.com/itsonlycode/gosecret/pkg/ctxutil"
	"github.com/itsonlycode/gosecret/pkg/pwgen"
	"github.com/itsonlycode/gosecret/tests/gptest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnvLeafHappyPath(t *testing.T) {
	u := gptest.NewUnitTester(t)
	defer u.Remove()

	ctx := context.Background()
	ctx = ctxutil.WithAlwaysYes(ctx, true)
	ctx = ctxutil.WithTerminal(ctx, false)
	act, err := newMock(ctx, u)
	require.NoError(t, err)
	require.NotNil(t, act)

	buf := &bytes.Buffer{}
	out.Stdout = buf
	out.Stderr = buf
	stdout = buf
	defer func() {
		out.Stdout = os.Stdout
		out.Stderr = os.Stderr
		stdout = os.Stdout
	}()

	// Command-line would be: "gosecret env foo env", where "foo" is an existing
	// secret with value "secret". We expect to see the key/value in the output
	// of the /usr/bin/env utility in the form "BAZ=secret".
	pw := pwgen.GeneratePassword(24, false)
	assert.NoError(t, act.insertStdin(ctx, "baz", []byte(pw), false))
	buf.Reset()

	assert.NoError(t, act.Env(gptest.CliCtx(ctx, t, "baz", "env")))
	assert.Contains(t, buf.String(), fmt.Sprintf("BAZ=%s\n", pw))
}

func TestEnvSecretNotFound(t *testing.T) {
	u := gptest.NewUnitTester(t)
	defer u.Remove()

	ctx := context.Background()
	ctx = ctxutil.WithAlwaysYes(ctx, true)
	ctx = ctxutil.WithTerminal(ctx, false)
	act, err := newMock(ctx, u)
	require.NoError(t, err)
	require.NotNil(t, act)

	// Command-line would be: "gosecret env non-existing true".
	assert.EqualError(t, act.Env(gptest.CliCtx(ctx, t, "non-existing", "true")),
		"Secret non-existing not found")
}

func TestEnvProgramNotFound(t *testing.T) {
	u := gptest.NewUnitTester(t)
	defer u.Remove()

	ctx := context.Background()
	ctx = ctxutil.WithAlwaysYes(ctx, true)
	ctx = ctxutil.WithTerminal(ctx, false)
	act, err := newMock(ctx, u)
	require.NoError(t, err)
	require.NotNil(t, act)

	wanted := "exec: \"non-existing\": executable file not found in "
	if runtime.GOOS == "windows" {
		wanted += "%PATH%"
	} else {
		wanted += "$PATH"
	}

	// Command-line would be: "gosecret env foo non-existing".
	assert.EqualError(t, act.Env(gptest.CliCtx(ctx, t, "foo", "non-existing")),
		wanted)
}

// Crash regression
func TestEnvProgramNotSpecified(t *testing.T) {
	u := gptest.NewUnitTester(t)
	defer u.Remove()

	ctx := context.Background()
	ctx = ctxutil.WithAlwaysYes(ctx, true)
	ctx = ctxutil.WithTerminal(ctx, false)
	act, err := newMock(ctx, u)
	require.NoError(t, err)
	require.NotNil(t, act)

	// Command-line would be: "gosecret env foo".
	assert.EqualError(t, act.Env(gptest.CliCtx(ctx, t, "foo")),
		"Missing subcommand to execute")
}
