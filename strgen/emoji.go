package strgen

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

var emojis = []rune{'😀', '😁', '😂', '🤣', '😃', '😄', '😅', '😆', '😉', '😊',
	'😋', '😎', '😍', '😘', '🥰', '😗', '😙', '😚', '🙂', '🤗', '🤩', '🤔', '🤨', '😐',
	'😑', '😶', '🙄', '😏', '😣', '😥', '😮', '🤐', '😯', '😪', '😫', '😴', '😌', '😛',
	'😜', '😝', '🤤', '😒', '😓', '😔', '😕', '🙃', '🤑', '😲', '🙁', '😖', '😞', '😟',
	'😤', '😢', '😭', '😦', '😧', '😨', '😩', '🤯', '😬', '😰', '😱', '🥵', '🥶', '😳',
	'🤪', '😵', '😡', '😠', '🤬', '😷', '🤒', '🤕', '🤢', '🤮', '🤧', '😇', '🤠', '🤡',
	'🥳', '🥴', '🥺', '🤥', '🤫', '🤭', '🧐', '🤓', '😈', '👿', '👹', '👺', '💀', '👻',
	'👽', '🤖', '💩', '😺', '😸', '😹', '😻', '😼', '😽', '🙀', '😿', '😾'}

// Emojifier is a Generator that produces emojis.
type Emojifier struct {
	Mode  string
	Index int
}

// NewEmojifier makes a new Emojifier
func NewEmojifier(mode string) (*Emojifier, error) {
	switch mode {
	case "":
		mode = "single"
	case "single", "staircase":
	default:
		return nil, fmt.Errorf("strgen: illegal Emojifier mode '%s'", mode)
	}

	return &Emojifier{Mode: mode}, nil
}

// Generate generates Emojis based on e.Mode.
func (e *Emojifier) Generate() string {
	var s string

	switch e.Mode {
	case "staircase":
		s = e.RandomN(e.Index + 1)
	default:
		s = e.RandomN(1)
	}

	e.Index++
	return s
}

// HasMore indicates whether an Emojifier has more emojis to generate (always
// true).
func (e *Emojifier) HasMore() bool {
	return true
}

// RandomN returns a string consisting of n pseudorandom emojis.
func (e *Emojifier) RandomN(n int) string {
	rand.NewSource(time.Now().UnixNano())

	builder := new(strings.Builder)
	for i := 0; i < n; i++ {
		eind := rand.Intn(len(emojis))
		if _, err := builder.WriteRune(emojis[eind]); err != nil {
			panic(err)
		}
	}

	return builder.String()
}
