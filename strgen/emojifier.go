package strgen

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/stevenxie/begone/emoji"
	ess "github.com/unixpickle/essentials"
)

// Emojifier is a Generator that produces emojis.
type Emojifier struct {
	Mode string // defines how emojis are generated

	index int    // keeps track of how many iterations have passed
	src   []rune // a set of emojis available for use

	rng *rand.Rand // random number source
}

// NewEmojifier makes a new Emojifier
func NewEmojifier(mode string) (*Emojifier, error) {
	// Default to mode "single".
	if mode == "" {
		mode = "single"
	}

	switch mode {
	case "single", "staircase", "unique":
	default:
		return nil, fmt.Errorf("strgen: unknown mode '%s'", mode)
	}

	src := rand.NewSource(time.Now().UnixNano())
	return &Emojifier{
		Mode: mode,
		src:  emoji.GetEmojis(),
		rng:  rand.New(src),
	}, nil
}

// Generate generates emojis, using a pattern based on e.Mode.
func (e *Emojifier) Generate() (string, error) {
	var (
		msg string
		err error
	)

	switch e.Mode {
	case "staircase":
		if msg, err = e.RandomN(e.index + 1); err != nil {
			return "", err
		}

	case "unique":
		if e.index < len(e.src) {
			msg = string(e.src[e.index])
		}

	default:
		if msg, err = e.RandomN(1); err != nil {
			return "", err
		}
	}

	e.index++
	return msg, nil
}

// HasMore indicates whether the Emojifier has more emojis to generate.
func (e *Emojifier) HasMore() bool {
	switch e.Mode {
	case "unique":
		return e.index < len(e.src)
	default:
		return true
	}
}

// RandomN returns a string consisting of n pseudorandom emojis.
func (e *Emojifier) RandomN(n int) (string, error) {
	builder := new(strings.Builder)
	for i := 0; i < n; i++ {
		emojiIndex := e.rng.Intn(len(e.src))
		if _, err := builder.WriteRune(e.src[emojiIndex]); err != nil {
			return "", ess.AddCtx("strgen: building emoji string", err)
		}
	}
	return builder.String(), nil
}
