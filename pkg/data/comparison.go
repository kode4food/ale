package data

// Comparison represents the result of an equality comparison
type Comparison int

// Comparison results
const (
	LessThan Comparison = iota - 1
	EqualTo
	GreaterThan
	Incomparable
)
