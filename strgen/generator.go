package strgen

// Generator generates strings.
type Generator interface {
	// Generate produces a string.
	Generate() string

	// HasMore returns true if Generator has more output to produce, and false
	// otherwise.
	HasMore() bool
}
