package config

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/itsonlycode/gosecret/pkg/appdir"
	"github.com/itsonlycode/gosecret/pkg/debug"
	"github.com/itsonlycode/gosecret/pkg/fsutil"

	homedir "github.com/mitchellh/go-homedir"
)

// Homedir returns the users home dir or an empty string if the lookup fails
func Homedir() string {
	if hd := os.Getenv("GOPASS_HOMEDIR"); hd != "" {
		return hd
	}
	hd, err := homedir.Dir()
	if err != nil {
		debug.Log("Failed to get homedir: %s\n", err)
		return ""
	}
	return hd
}

// configLocation returns the location of the config file
// (a YAML file that contains values such as the path to the password store)
func configLocation() string {
	// First, check for the "GOPASS_CONFIG" environment variable
	if cf := os.Getenv("GOPASS_CONFIG"); cf != "" {
		return cf
	}

	// Second, check for the "XDG_CONFIG_HOME" environment variable
	// (which is part of the XDG Base Directory Specification for Linux and
	// other Unix-like operating sytstems)
	return filepath.Join(appdir.UserConfig(), "config.yml")
}

// configLocations returns the possible locations of gosecret config files,
// in decreasing priority
func configLocations() []string {
	l := []string{}
	if cf := os.Getenv("GOPASS_CONFIG"); cf != "" {
		l = append(l, cf)
	}
	l = append(l, filepath.Join(appdir.UserConfig(), "config.yml"))
	l = append(l, filepath.Join(Homedir(), ".config", "gosecret", "config.yml"))
	l = append(l, filepath.Join(Homedir(), ".gosecret.yml"))
	return l
}

// PwStoreDir reads the password store dir from the environment
// or returns the default location if the env is not set
func PwStoreDir(mount string) string {
	if mount != "" {
		cleanName := strings.Replace(mount, string(filepath.Separator), "-", -1)
		return fsutil.CleanPath(filepath.Join(appdir.UserData(), "stores", cleanName))
	}
	// PASSWORD_STORE_DIR support is discouraged
	if d := os.Getenv("PASSWORD_STORE_DIR"); d != "" {
		return fsutil.CleanPath(d)
	}
	if ld := filepath.Join(appdir.UserHome(), ".password-store"); fsutil.IsDir(ld) {
		debug.Log("re-using existing legacy dir for root store: %s", ld)
		return ld
	}
	return fsutil.CleanPath(filepath.Join(appdir.UserData(), "stores", "root"))
}

// Directory returns the configuration directory for the gosecret config file
func Directory() string {
	return filepath.Dir(configLocation())
}
