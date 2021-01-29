package data

import "bytes"

// Vector is a fixed-length array of Values
type Vector []Value

// EmptyVector represents an empty Vector
var EmptyVector = Vector{}

// NewVector creates a new Vector instance
func NewVector(v ...Value) Vector {
	return v
}

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

// IsEmpty returns whether this sequence is empty
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

// Car returns the first element of a Pair
func (v Vector) Car() Value {
	return SequenceCar(v)
}

// Cdr returns the second element of a Pair
func (v Vector) Cdr() Value {
	return SequenceCdr(v)
}

// Prepend inserts an element at the beginning of the Vector
func (v Vector) Prepend(e Value) Sequence {
	return append(Vector{e}, v...)
}

// Append appends elements to the end of the Vector
func (v Vector) Append(e Value) Sequence {
	return append(v, e)
}

// Reverse returns a reversed copy of this Vector
func (v Vector) Reverse() Sequence {
	vl := len(v)
	if vl <= 1 {
		return v
	}
	res := make(Vector, vl)
	for i, j := 0, vl-1; j >= 0; i, j = i+1, j-1 {
		res[i] = v[j]
	}
	return res
}

// Call turns Vector into a Function
func (v Vector) Call(args ...Value) Value {
	return indexedCall(v, args)
}

// Convention returns the Function's calling convention
func (v Vector) Convention() Convention {
	return ApplicativeCall
}

// CheckArity performs a compile-time arity check for the Function
func (v Vector) CheckArity(argCount int) error {
	return checkRangedArity(1, 2, argCount)
}

// Equal compares this Vector to another for equality
func (v Vector) Equal(r Value) bool {
	if r, ok := r.(Vector); ok {
		if len(v) != len(r) {
			return false
		}
		for i, elem := range r {
			if !v[i].Equal(elem) {
				return false
			}
		}
		return true
	}
	return false
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
