package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// Remove removes the config file, if it exists.
// It returns a list of the files that were removed.
func Remove() (removed []string, err error) {
	path, err := Filepath()
	if err != nil {
		return nil, err
	}

	if err = os.Remove(path); err != nil {
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("config: failed to remove file at '%s': %v", path,
				err)
		}
	} else {
		removed = []string{path}
	}

	if runtime.GOOS == "windows" {
		dirPath := filepath.Dir(path)
		if err = os.Remove(dirPath); err != nil {
			if !os.IsNotExist(err) {
				return nil, fmt.Errorf("config: failed to remove config directory at "+
					"'%s': %v", path, err)
			}
		} else {
			removed = append(removed, dirPath)
		}
	}

	return removed, nil
}
