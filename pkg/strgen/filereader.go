package strgen

import (
	"bufio"
	"errors"
	"os"
	"strings"

	ess "github.com/unixpickle/essentials"
)

// FileReader is an mbot.StringGenerator that reads from a file.
type FileReader struct {
	Scanner *bufio.Scanner
	File    *os.File
	Path    string
}

// NewFileReader returns a new FileReader.
func NewFileReader(fpath string) (*FileReader, error) {
	info, err := os.Stat(fpath)
	if err != nil {
		return nil, ess.AddCtx("strgen: getting file info", err)
	}

	if info.IsDir() {
		return nil, errors.New("strgen: file can't be a directory")
	}
	return &FileReader{Path: fpath}, nil
}

// Generate returns an unread line from fr.File.
func (fr *FileReader) Generate() (string, error) {
	if fr.Scanner == nil {
		var err error
		fr.File, err = os.Open(fr.Path)
		if err != nil {
			return "", ess.AddCtx("strgen: opening file", err)
		}

		fr.Scanner = bufio.NewScanner(fr.File)
		fr.Scanner.Scan()
	}

	if err := fr.Scanner.Err(); err != nil {
		return "", ess.AddCtx("strgen: error from scanner", err)
	}
	return fr.Scanner.Text(), nil
}

// HasMore returns true if fr.Scanner has more lines to scan, if it is
// uninitialized, or if it has an error.
func (fr *FileReader) HasMore() bool {
	if fr.Scanner == nil {
		return true
	}

	var text string
	for text == "" {
		if !fr.Scanner.Scan() { // scan did not advance to next token
			if err := fr.File.Close(); err != nil {
				panic(err)
			}
			return fr.Scanner.Err() != nil
		}
		text = strings.TrimSpace(fr.Scanner.Text())
	}

	return true
}
