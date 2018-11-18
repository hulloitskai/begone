package loader

import (
	"archive/zip"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// DownloadTo downloads the OS-specific ChromeDriver to dst (a filepath).
//
// The installation directory must already exist.
func (l *Loader) DownloadTo(dst string) error {
	// Guard against bad destination.
	dstDir := filepath.Dir(dst)
	if _, err := os.Stat(dstDir); os.IsNotExist(err) {
		return fmt.Errorf("loader: installation directory '%s' does not exist",
			dstDir)
	}

	// Generate filename.
	var (
		dlname = fmt.Sprintf("chromedriver-%s.zip", timeString()) // is unique
		dlpath = filepath.Join(os.TempDir(), dlname)
	)

	// Assert that dlpath is absolute (for context during errors).
	if !filepath.IsAbs(dlpath) {
		return fmt.Errorf("loader: expected generated dlpath to be an absolute "+
			"path, but got: %s", dlpath)
	}

	dlfile, err := os.Create(dlpath)
	if err != nil {
		return fmt.Errorf("loader: failed to create dlfile '%s': %v", dlpath, err)
	}
	defer dlfile.Close()

	res, err := http.Get(l.DriverURL)
	if err != nil {
		return fmt.Errorf("loader: failed to accessing driver URL (%s): %v",
			l.DriverURL, err)
	}
	defer res.Body.Close()

	// Write body to dlfile and hash.
	hash := sha256.New()
	if _, err = io.Copy(io.MultiWriter(dlfile, hash), res.Body); err != nil {
		return fmt.Errorf("loader: failed to write archive: %v", err)
	}

	if err = res.Body.Close(); err != nil {
		return fmt.Errorf("loader: failed to close response body: %v", err)
	}

	// Verify SHA256 checksum.
	if sum := hex.EncodeToString(hash.Sum(nil)); sum != l.DriverHash {
		return fmt.Errorf("loader: unexpected MD5 checksum of downloaded archive "+
			"at '%s'; expected '%s', instead got '%s'", dlpath, l.DriverHash, sum)
	}

	// Unzip dlfile to dst.
	info, err := dlfile.Stat()
	if err != nil {
		return fmt.Errorf("loader: failed to read dlfile info: %v", err)
	}

	reader, err := zip.NewReader(dlfile, info.Size())
	if err != nil {
		return fmt.Errorf("loader: failed to create zip reader: %v", err)
	}

	if len(reader.File) != 1 {
		return fmt.Errorf("loader: expected downloaded archive to contain 1 "+
			"file, instead got %d", len(reader.File))
	}

	ifile, err := reader.File[0].Open()
	if err != nil {
		return fmt.Errorf("loader: failed to open archive file for reading: %v",
			err)
	}
	defer ifile.Close()

	ofile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("loader: failed to create ofile '%s': %v", dst, err)
	}
	defer ofile.Close()

	if _, err = io.Copy(ofile, ifile); err != nil {
		return fmt.Errorf("loader: failed to write from archive to ofile: %v", err)
	}

	if err = dlfile.Close(); err != nil {
		return fmt.Errorf("loader: failed to close dlfile: %v", err)
	}
	if err = os.Remove(dlpath); err != nil {
		return fmt.Errorf("loader: could not remove dlfile: %v", err)
	}

	// Set ofile permissions (make executable).
	if err = ofile.Chmod(0755); err != nil {
		return fmt.Errorf("loader: failed to modify ofile: %v", err)
	}

	// Close files.
	if err = ifile.Close(); err != nil {
		return fmt.Errorf("loader: could not close ifile: %v", err)
	}
	if err = ofile.Close(); err != nil {
		return fmt.Errorf("loader: could not close ofile: %v", err)
	}

	return nil
}
