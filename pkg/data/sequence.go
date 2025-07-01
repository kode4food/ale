package data

type (
	// Sequence interfaces expose a lazily resolved sequence
	Sequence interface {
		Pair
		Split() (Value, Sequence, bool)
		IsEmpty() bool
	}

	// Appender is a Sequence that can be appended to
	Appender interface {
		Sequence
		Append(Value) Sequence
	}

	// Mapper is a Sequence that provides a mutable Mapped interface
	Mapper interface {
		Sequence
		Mapped
		Put(Pair) Sequence
		Remove(Value) (Value, Sequence, bool)
	}

	// Prepender is a Sequence that can be prepended to
	Prepender interface {
		Sequence
		Prepend(Value) Sequence
	}

	// Reverser is a Sequence than can be reversed
	Reverser interface {
		Sequence
		Reverse() Sequence
	}

	// RandomAccess provides a Sequence that supports random access
	RandomAccess interface {
		Indexed
		Counted
	}

	// Counted is a Sequence that returns a count of its items
	Counted interface {
		Sequence
		Count() int
	}

	// Indexed is a Sequence that has indexed elements
	Indexed interface {
		Sequence
		ElementAt(int) (Value, bool)
	}
)
