package data

type (
	// Comparison represents the result of an equality comparison
	Comparison int

	// Comparer is an interface for a Value capable of comparing
	Comparer interface {
		Compare(Comparer) Comparison
	}
)

// Comparison results
const (
	LessThan Comparison = iota - 1
	EqualTo
	GreaterThan
	Incomparable
)
