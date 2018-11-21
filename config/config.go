package config

import (
	"os"
	"os/user"
	"path/filepath"
	"runtime"

	ess "github.com/unixpickle/essentials"
)

// Config file path info.
const (
	FileName = ".begone.json"
	FileDir  = "me.stevenxie.begone"
)

// Config contains persistent program options for begone.
type Config struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Filepath returns the path to the config file.
//
// The config file will reside in the active user's home directory on macOS and
// Linux, and in the %APPDATA% directory on Windows.
func Filepath() (string, error) {
	u, err := user.Current()
	if err != nil {
		return "", ess.AddCtx("config: getting current user", err)
	}

	switch runtime.GOOS {
	case "windows":
		roamingDir, ok := os.LookupEnv("APPDATA")
		if !ok {
			roamingDir = filepath.Join(u.HomeDir, "AppData", "Roaming")
		}

		// Determine namespace directory in %APPDATA%, and make it if it does not
		// yet exist.
		nsDir := filepath.Join(roamingDir)
		if _, err = os.Stat(nsDir); os.IsNotExist(err) {
			if err = os.Mkdir(nsDir, 0755); err != nil {
				return "", ess.AddCtx("config: creating config directory", err)
			}
		}
		return filepath.Join(nsDir, FileName), nil

	case "linux":
		cfgHome, ok := os.LookupEnv("XDG_CONFIG_HOME")
		if !ok {
			cfgHome = u.HomeDir
		}
		return filepath.Join(cfgHome, FileName), nil

	default:
		return filepath.Join(u.HomeDir, FileName), nil
	}
}
