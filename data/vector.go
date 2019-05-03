package data

import "bytes"

// Vector is a fixed-length array of Values
type Vector []Value

// EmptyVector is an empty Vector
var EmptyVector = Vector{}

// Count returns the number of elements in the Vector
func (v Vector) Count() int {
	return len(v)
}

// ElementAt returns a specific element of the Vector
func (v Vector) ElementAt(index int) (Value, bool) {
	if index >= 0 && index < len(v) {
		return v[index], true
	}
	return Nil, false
}

// First returns the first element of the Vector
func (v Vector) First() Value {
	if len(v) > 0 {
		return v[0]
	}
	return Nil
}

// Rest returns the elements of the Vector that follow the first
func (v Vector) Rest() Sequence {
	if len(v) > 1 {
		return v[1:]
	}
	return EmptyVector
}

// IsEmpty returns whether or not this sequence is empty
func (v Vector) IsEmpty() bool {
	return len(v) == 0
}

// Split breaks the Vector into its components (first, rest, ok)
func (v Vector) Split() (Value, Sequence, bool) {
	lv := len(v)
	if lv > 1 {
		return v[0], v[1:], true
	} else if lv == 1 {
		return v[0], EmptyVector, true
	}
	return Nil, EmptyVector, false
}

// Prepend inserts an element at the beginning of the Vector
func (v Vector) Prepend(p Value) Sequence {
	return append(Vector{p}, v...)
}

// Append appends an element to the end of the Vector
func (v Vector) Append(args ...Value) Sequence {
	return append(v, args...)
}

// Reverse returns a reversed copy of this Vector
func (v Vector) Reverse() Sequence {
	res := make(Vector, len(v))
	for i, j := 0, len(v)-1; j >= 0; i, j = i+1, j-1 {
		res[i] = v[j]
	}
	return res
}

// Caller turns Vector into a callable type
func (v Vector) Caller() Call {
	return makeIndexedCall(v)
}

// String converts this Vector to a string
func (v Vector) String() string {
	var b bytes.Buffer
	l := len(v)

	b.WriteString("[")
	for i := 0; i < l; i++ {
		if i > 0 {
			b.WriteString(" ")
		}
		b.WriteString(MaybeQuoteString(v[i]))
	}
	b.WriteString("]")
	return b.String()
}
