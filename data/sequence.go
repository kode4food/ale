package data

import "bytes"

type (
	// Sequence interfaces expose a lazily resolved sequence of Vector
	Sequence interface {
		Value
		First() Value
		Rest() Sequence
		Split() (Value, Sequence, bool)
		Prepend(Value) Sequence
		IsEmpty() bool
	}

	// Appender is a Sequence that can be appended to
	Appender interface {
		Sequence
		Append(...Value) Sequence
	}

	// Reverser is a Sequence than can be reversed
	Reverser interface {
		Sequence
		Reverse() Sequence
	}

	// IndexedSequence is a Sequence that provides an Indexed interface
	IndexedSequence interface {
		Sequence
		Indexed
	}

	// SizedSequence is a Sequence that provides a Sized interface
	SizedSequence interface {
		Sequence
		Sized
	}

	// RandomAccessSequence provides Indexed and Sized Sequence interfaces
	RandomAccessSequence interface {
		Sequence
		Indexed
		Sized
	}

	// MappedSequence is a Sequence that provides a Mapped interface
	MappedSequence interface {
		Sequence
		Mapped
	}
)

// MakeSequenceStr converts a Sequence to a String
func MakeSequenceStr(s Sequence) string {
	f, r, ok := s.Split()
	if !ok {
		return "()"
	}

	var b bytes.Buffer
	b.WriteString("(")
	b.WriteString(MaybeQuoteString(f))
	for f, r, ok = r.Split(); ok; f, r, ok = r.Split() {
		b.WriteString(" ")
		b.WriteString(MaybeQuoteString(f))
	}
	b.WriteString(")")
	return b.String()
}

// Size will return the size from a Sized Sequence or explode
func Size(v Value) int {
	return v.(SizedSequence).Size()
}

// Last returns the final element of a Sequence, possibly by scanning
func Last(s Sequence) (Value, bool) {
	if s.IsEmpty() {
		return Nil, false
	}

	if i, ok := s.(RandomAccessSequence); ok {
		return i.ElementAt(i.Size() - 1)
	}

	var res Value
	var lok bool
	for f, s, ok := s.Split(); ok; f, s, ok = s.Split() {
		res = f
		lok = ok
	}
	return res, lok
}

func makeIndexedCall(s IndexedSequence) Call {
	return func(args ...Value) Value {
		idx := args[0].(Integer)
		res, ok := s.ElementAt(int(idx))
		if !ok && len(args) > 1 {
			return args[1]
		}
		return res
	}
}

func makeMappedCall(m Mapped) Call {
	return func(args ...Value) Value {
		res, ok := m.Get(args[0])
		if !ok && len(args) > 1 {
			return args[1]
		}
		return res
	}
}
