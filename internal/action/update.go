package action

import (
	"fmt"
	"runtime"

	"github.com/itsonlycode/gosecret/internal/out"
	"github.com/itsonlycode/gosecret/internal/updater"
	"github.com/itsonlycode/gosecret/pkg/ctxutil"

	"github.com/urfave/cli/v2"
)

// Update will start the interactive update assistant
func (s *Action) Update(c *cli.Context) error {
	s.rem.Reset("update")

	ctx := ctxutil.WithGlobalFlags(c)

	if s.version.String() == "0.0.0+HEAD" {
		out.Errorf(ctx, "Can not check version against HEAD")
		return nil
	}

	if runtime.GOOS == "windows" {
		return fmt.Errorf("gosecret update is not supported on windows (#1722)")
	}

	out.Printf(ctx, "⚒ Checking for available updates ...")
	if err := updater.Update(ctx, s.version); err != nil {
		return ExitError(ExitUnknown, err, "Failed to update gosecret: %s", err)
	}

	out.OKf(ctx, "gosecret is up to date")
	return nil
}
