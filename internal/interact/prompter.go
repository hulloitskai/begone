package interact

import (
	"io"
	"os"
	"runtime"
)

// Mode is the operating mode of a Prompter. It is one of:
//   - Standard
//   - Reduced
type Mode uint8

const (
	// Standard is the default operating mode, and enables usage of emojis,
	// spinners, and other 'fancy' terminal display elements.
	Standard Mode = iota

	// Reduced is a reduced-features operating mode that disables emojis,
	// spinners and other 'fancy' terminal display elements.
	Reduced
)

// Prompter produces command line prompts and reads user responses.
type Prompter struct {
	Out, Err io.Writer
	In       io.Reader
	Mode
}

// NewPrompter returns a Prompter with the default configuration.
//
// The resulting Prompter reads from os.Stdin, writes to os.Stdout, and selects
// a mode based on runtime.GOOS (reduced if running on Windows, Standard
// otherwise).
func NewPrompter() *Prompter {
	return NewPrompterWith(os.Stdout, os.Stderr)
}

// NewPrompterWith returns a Prompter that writes output to out, and errors to
// err.
//
// It selects a mode based on runtime.GOOS (reduced if running on Windows,
// Standard
// otherwise).
func NewPrompterWith(out, err io.Writer) *Prompter {
	var mode Mode
	if runtime.GOOS == "windows" {
		mode = Reduced
	}

	return &Prompter{
		In:  os.Stdin,
		Out: out, Err: err,
		Mode: mode,
	}
}
