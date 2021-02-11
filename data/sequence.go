package data

import "bytes"

type (
	// Sequence interfaces expose a lazily resolved sequence of Vector
	Sequence interface {
		Value
		First() Value
		Rest() Sequence
		Split() (Value, Sequence, bool)
		IsEmpty() bool
	}

	// AppenderSequence is a Sequence that acts as an Appender
	AppenderSequence interface {
		Sequence
		Appender
	}

	// CountedSequence is a Sequence that provides a Counted interface
	CountedSequence interface {
		Sequence
		Counted
	}

	// IndexedSequence is a Sequence that provides an Indexed interface
	IndexedSequence interface {
		Sequence
		Indexed
	}

	// MappedSequence is a Sequence that provides a Mapped interface
	MappedSequence interface {
		Sequence
		Mapped
	}

	// PrependerSequence is a Sequence that acts as a Prepender
	PrependerSequence interface {
		Sequence
		Prepender
	}

	// RandomAccessSequence provides a RandomAccess Sequence interface
	RandomAccessSequence interface {
		Sequence
		RandomAccess
	}

	// ReverserSequence is a Sequence than acts as a Reverser
	ReverserSequence interface {
		Sequence
		Reverser
	}

	// ValuerSequence is a Sequence that provides a Valuer interface
	ValuerSequence interface {
		Sequence
		Valuer
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

// Last returns the final element of a Sequence, possibly by scanning
func Last(s Sequence) (Value, bool) {
	if s.IsEmpty() {
		return Nil, false
	}

	if i, ok := s.(RandomAccessSequence); ok {
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

func indexedCall(s IndexedSequence, args []Value) Value {
	idx := args[0].(Integer)
	res, ok := s.ElementAt(int(idx))
	if !ok && len(args) > 1 {
		return args[1]
	}
	return res
}
