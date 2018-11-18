package loader

import (
	"fmt"
	"os"
	"path/filepath"
)

// Clean removes the installed driver (if it exists) along with the namespace
// folder, and also attempts to remove temporary artifacts resulting from
// failed download attempts.
//
// It returns a list of artifacts that were removed during the cleanup process.
func (l *Loader) Clean(removeDriver bool) (removed []string, err error) {
	// Locate temporary download artifacts.
	pat := filepath.Join(os.TempDir(), "chromedriver-*")
	paths, err := filepath.Glob(pat)
	if err != nil {
		return nil, fmt.Errorf("loader: error while matching files against glob "+
			"pattern: %v", err)
	}

	// Append installed driver.
	if removeDriver {
		if _, err = os.Stat(l.DriverPath); err != nil {
			if !os.IsNotExist(err) {
				return nil, fmt.Errorf("loader: failed to read driver file info: %v", err)
			}
		} else {
			paths = append(paths, l.DriverPath)
			paths = append(paths, filepath.Dir(l.DriverPath))
		}
	}

	// Remove selected files.
	for _, path := range paths {
		if err = os.Remove(path); err != nil {
			return removed, fmt.Errorf("loader: failed to remove file '%s': %v", path,
				err)
		}

		removed = append(removed, path) // keep track of successful removes
	}

	return removed, nil
}
