package data

import "github.com/kode4food/ale"

type (
	// Sequence interfaces expose a lazily resolved sequence
	Sequence interface {
		Pair

		// Split returns the split form (First and Rest) of the Sequence
		Split() (ale.Value, Sequence, bool)

		// IsEmpty returns whether this sequence is empty
		IsEmpty() bool
	}

	// Appender is a Sequence that can be appended to
	Appender interface {
		Sequence

		// Append creates a new Sequence with the given value appended to it
		Append(ale.Value) Sequence
	}

	// Mapper is a Sequence that provides a mutable Mapped interface
	Mapper interface {
		Sequence
		Mapped

		// Put returns a new Sequence with the Pair associated in the Mapper
		Put(Pair) Sequence

		// Remove returns a new Sequence with the Value removed from the Mapper
		Remove(ale.Value) (ale.Value, Sequence, bool)
	}

	// Prepender is a Sequence that can be prepended to
	Prepender interface {
		Sequence

		// Prepend creates a new Sequence with the given value prepended to it
		Prepend(ale.Value) Sequence
	}

	// Reverser is a Sequence that can be reversed
	Reverser interface {
		Sequence

		// Reverse creates a new Sequence with the elements in reverse order
		Reverse() Sequence
	}

	// Counted is a Sequence that returns a count of its items
	Counted interface {
		Sequence

		// Count returns the number of elements in this Sequence
		Count() int
	}

	// Indexed is a Sequence that has indexed elements
	Indexed interface {
		Counted

		// ElementAt returns the element at the specified index
		ElementAt(int) (ale.Value, bool)
	}
)
