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
		Sequence
		Indexed
		Counted
	}

	// IndexedSequence is a Sequence that has indexed elements
	IndexedSequence interface {
		Sequence
		Indexed
	}

	// CountedSequence is a Sequence that returns a count of its items
	CountedSequence interface {
		Sequence
		Counted
	}
)

// Last returns the final element of a Sequence, possibly by scanning
func Last(s Sequence) (Value, bool) {
	if s.IsEmpty() {
		return Null, false
	}

	if i, ok := s.(RandomAccess); ok {
		return i.ElementAt(i.Count() - 1)
	}

	var res Value
	var lok bool
	for f, s, ok := s.Split(); ok; f, s, ok = s.Split() {
		res = f
		lok = ok
	}
	return res, lok
}

func indexedCall(s IndexedSequence, args Vector) Value {
	idx := args[0].(Integer)
	res, ok := s.ElementAt(int(idx))
	if !ok && len(args) > 1 {
		return args[1]
	}
	return res
}

func mappedCall(m Mapper, args Vector) Value {
	res, ok := m.Get(args[0])
	if !ok && len(args) > 1 {
		return args[1]
	}
	return res
}
