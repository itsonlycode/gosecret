package action

import (
	"io"
	"os"
	"path/filepath"

	"github.com/itsonlycode/gosecret/internal/config"
	"github.com/itsonlycode/gosecret/internal/reminder"
	"github.com/itsonlycode/gosecret/internal/store/root"
	"github.com/itsonlycode/gosecret/pkg/debug"

	"github.com/blang/semver/v4"
)

var (
	stdin  io.Reader = os.Stdin
	stdout io.Writer = os.Stdout
)

// Action knows everything to run gosecret CLI actions
type Action struct {
	Name    string
	Store   *root.Store
	cfg     *config.Config
	version semver.Version
	rem     *reminder.Store
}

// New returns a new Action wrapper
func New(cfg *config.Config, sv semver.Version) (*Action, error) {
	return newAction(cfg, sv, true)
}

func newAction(cfg *config.Config, sv semver.Version, remind bool) (*Action, error) {
	name := "gosecret"
	if len(os.Args) > 0 {
		name = filepath.Base(os.Args[0])
	}

	act := &Action{
		Name:    name,
		cfg:     cfg,
		version: sv,
		Store:   root.New(cfg),
	}

	if remind {
		r, err := reminder.New()
		if err != nil {
			debug.Log("failed to init reminder: %s", err)
		} else {
			// only populate the reminder variable on success, the implementation
			// can handle being called on a nil pointer
			act.rem = r
		}
	}

	return act, nil
}

// String implement fmt.Stringer
func (s *Action) String() string {
	return s.Store.String()
}
