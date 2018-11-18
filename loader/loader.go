package loader

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
)

// Loader is capable of installing a ChromeDriver for Selenium into the host OS.
type Loader struct {
	DriverURL  string
	DriverHash string
	DriverPath string
	namespace  string
	TempDir    string
}

// NewLoader returns a Loader, configured for the host OS.
func NewLoader(namespace string) (*Loader, error) {
	// Get user info.
	u, err := user.Current()
	if err != nil {
		return nil, fmt.Errorf("loader: couldn't get current user info: %v", err)
	}

	// Create platform-specific loaders.
	switch runtime.GOOS {
	case "darwin":
		// Determine driver installation path.
		path := filepath.Join(u.HomeDir, "Library", "Caches", namespace, DriverName)

		return &Loader{
			DriverURL:  MacDriverURL,
			DriverHash: MacDriverHash,
			DriverPath: path,
		}, nil

	case "linux":
		cachePath, ok := os.LookupEnv("XDG_CACHE_HOME")
		if !ok {
			cachePath = filepath.Join(u.HomeDir, ".cache")
		}

		path := filepath.Join(cachePath, namespace, DriverName)
		return &Loader{
			DriverURL:  LinuxDriverURL,
			DriverHash: LinuxDriverHash,
			DriverPath: path,
		}, nil

	case "windows":
		cachePath, ok := os.LookupEnv("LOCALAPPDATA")
		if !ok {
			cachePath = filepath.Join(u.HomeDir, "AppData", "Local")
		}

		path := filepath.Join(cachePath, namespace, WinDriverName)
		return &Loader{
			DriverURL:  WinDriverURL,
			DriverHash: WinDriverHash,
			DriverPath: path,
		}, nil
	}

	return nil, fmt.Errorf("loader: no Loader available for this OS (%v)",
		runtime.GOOS)
}
