package loader

import (
	"fmt"
	"os"
	"path/filepath"
)

// Ensure ensures that the Selenium ChromeDriver is installed in the
// expected directory.
func (l *Loader) Ensure() error {
	// Check the driver, download it if it does not exist.
	info, err := os.Stat(l.DriverPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("loader: could not get driver file info: %v", err)
		}

		// It must be that the driver has not yet been installed. So, install it!
		// Create installation folder if it does not exist.
		driverDir := filepath.Dir(l.DriverPath)
		if _, err = os.Stat(driverDir); err != nil {
			if !os.IsNotExist(err) {
				return fmt.Errorf("loader: error getting driver directory info: %v",
					err)
			}

			if err = os.MkdirAll(driverDir, 0755); err != nil {
				return fmt.Errorf("loader: error making driver installation "+
					"directory (%s): %v", driverDir, err)
			}
		}

		if err = l.DownloadTo(l.DriverPath); err != nil {
			return err
		}

		if info, err = os.Stat(l.DriverPath); err != nil {
			return fmt.Errorf("loader: error getting driver file info after "+
				"successful installation: %v", err)
		}
	}

	// Verify that driver is excutable.
	if (info.Mode() & 0100) != 0100 {
		return fmt.Errorf("loader: driver does not appear to be executable (got "+
			"file mode: %o)", info.Mode())
	}

	return nil
}
