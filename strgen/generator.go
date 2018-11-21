package strgen

// Generator generates strings.
type Generator interface {
	// Generate returns a string.
	Generate() (string, error)

	// HasMore returns true if Generator has more output to produce, and false
	// otherwise.
	HasMore() bool
}
