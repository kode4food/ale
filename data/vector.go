package data

import "bytes"

type (
	// Vector is a fixed-length array of Values
	Vector interface {
		vector() // marker
		Sequence
		RandomAccess
		Prepender
		Appender
		Reverser
		Valuer
	}

	vector []Value
)

// EmptyVector represents an empty Vector
var EmptyVector = NewVector()

// NewVector creates a new Vector instance
func NewVector(v ...Value) Vector {
	return vector(v)
}

func (vector) vector() {}

func (v vector) Values() Values {
	return Values(v)
}

// Count returns the number of elements in the Vector
func (v vector) Count() int {
	return len(v)
}

// ElementAt returns a specific element of the Vector
func (v vector) ElementAt(index int) (Value, bool) {
	if index >= 0 && index < len(v) {
		return v[index], true
	}
	return Nil, false
}

// First returns the first element of the Vector
func (v vector) First() Value {
	if len(v) > 0 {
		return v[0]
	}
	return Nil
}

// Rest returns the elements of the Vector that follow the first
func (v vector) Rest() Sequence {
	if len(v) > 1 {
		return v[1:]
	}
	return EmptyVector
}

// IsEmpty returns whether this sequence is empty
func (v vector) IsEmpty() bool {
	return len(v) == 0
}

// Split breaks the Vector into its components (first, rest, ok)
func (v vector) Split() (Value, Sequence, bool) {
	lv := len(v)
	if lv > 1 {
		return v[0], v[1:], true
	} else if lv == 1 {
		return v[0], EmptyVector, true
	}
	return Nil, EmptyVector, false
}

// Car returns the first element of a Pair
func (v vector) Car() Value {
	return SequenceCar(v)
}

// Cdr returns the second element of a Pair
func (v vector) Cdr() Value {
	return SequenceCdr(v)
}

// Prepend inserts an element at the beginning of the Vector
func (v vector) Prepend(e Value) Sequence {
	return append(vector{e}, v...)
}

// Append appends elements to the end of the Vector
func (v vector) Append(e Value) Sequence {
	return append(v, e)
}

// Reverse returns a reversed copy of this Vector
func (v vector) Reverse() Sequence {
	vl := len(v)
	if vl <= 1 {
		return v
	}
	res := make(vector, vl)
	for i, j := 0, vl-1; j >= 0; i, j = i+1, j-1 {
		res[i] = v[j]
	}
	return res
}

// Call turns Vector into a Function
func (v vector) Call(args ...Value) Value {
	return indexedCall(v, args)
}

// Convention returns the Function's calling convention
func (v vector) Convention() Convention {
	return ApplicativeCall
}

// CheckArity performs a compile-time arity check for the Function
func (v vector) CheckArity(argCount int) error {
	return checkRangedArity(1, 2, argCount)
}

// Equal compares this Vector to another for equality
func (v vector) Equal(r Value) bool {
	if r, ok := r.(vector); ok {
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
func (v vector) String() string {
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
