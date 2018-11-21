package strgen

// Repeater is a Generator that simply returns the same string every time.
type Repeater struct {
	Msg string
}

// NewRepeater creates a Repeater from s.
func NewRepeater(s string) *Repeater {
	return &Repeater{Msg: s}
}

// HasMore indicates whether or not the repeater has more text to generate
// (always true).
func (r *Repeater) HasMore() bool {
	return true
}

// Generate returns the string to repeat (r.Msg).
func (r Repeater) Generate() (string, error) {
	return r.Msg, nil
}
