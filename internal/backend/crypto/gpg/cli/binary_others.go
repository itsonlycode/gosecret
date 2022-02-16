//go:build !windows
// +build !windows

package cli

import (
	"os/exec"

	"github.com/itsonlycode/gosecret/pkg/debug"
	"github.com/itsonlycode/gosecret/pkg/fsutil"
	"github.com/itsonlycode/pinentry/gpgconf"
)

func detectBinary(name string) (string, error) {
	// user supplied binaries take precedence
	if name != "" {
		return exec.LookPath(name)
	}
	// try to get the proper binary from gpgconf(1)
	p, err := gpgconf.Path("gpg")
	if err != nil || p == "" || !fsutil.IsFile(p) {
		debug.Log("gpgconf failed (%q), falling back to path lookup: %q", p, err)
		// otherwise fall back to the default and try
		// to look up "gpg"
		return exec.LookPath("gpg")
	}

	debug.Log("gpgconf returned %q for gpg", p)
	return p, nil
}
