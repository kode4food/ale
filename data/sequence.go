package data

import "github.com/kode4food/ale"

type (
	// Sequence interfaces expose a lazily resolved sequence
	Sequence interface {
		Pair
		Split() (ale.Value, Sequence, bool)
		IsEmpty() bool
	}

	// Appender is a Sequence that can be appended to
	Appender interface {
		Sequence
		Append(ale.Value) Sequence
	}

	// Mapper is a Sequence that provides a mutable Mapped interface
	Mapper interface {
		Sequence
		Mapped
		Put(Pair) Sequence
		Remove(ale.Value) (ale.Value, Sequence, bool)
	}

	// Prepender is a Sequence that can be prepended to
	Prepender interface {
		Sequence
		Prepend(ale.Value) Sequence
	}

	// Reverser is a Sequence than can be reversed
	Reverser interface {
		Sequence
		Reverse() Sequence
	}

	// Counted is a Sequence that returns a count of its items
	Counted interface {
		Sequence
		Count() int
	}

	// Indexed is a Sequence that has indexed elements
	Indexed interface {
		Counted
		ElementAt(int) (ale.Value, bool)
	}
)
