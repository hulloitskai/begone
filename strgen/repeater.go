package strgen

// Repeater is a Generator that simply returns the same string every time.
type Repeater string

// NewRepeater creates a Repeater from s.
func NewRepeater(s string) Repeater {
	return Repeater(s)
}

// HasMore returns true if Generator has more output to produce, and false
// otherwise.
func (r Repeater) HasMore() bool {
	return true
}

// Generate produces a string.
func (r Repeater) Generate() string {
	return string(r)
}
