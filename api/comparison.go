package api

type (
	// Comparison represents the result of a equality comparison
	Comparison int

	// Comparer is an interface for a Value capable of comparing
	Comparer interface {
		Compare(Comparer) Comparison
	}
)

//go:generate stringer -type=Comparison -linecomment
const (
	LessThan     Comparison = iota - 1 // Less Than
	EqualTo                            // Equal To
	GreaterThan                        // Greater Than
	Incomparable                       // Not Comparable
)
