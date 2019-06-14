package data

import (
	"bytes"
	"fmt"
)

// Associative is a Mapped Value that is implemented atop a Vector
type Associative []Vector

// Error messages
const (
	MapNotPaired = "map does not contain an even number of elements"
	ExpectedPair = "expected a key-value pair"
)

// EmptyAssociative represents an empty Associative
var EmptyAssociative = Associative{}

// NewAssociative instantiates a new Associative
func NewAssociative(v ...Value) Associative {
	l := len(v)
	if l%2 != 0 {
		panic(fmt.Errorf(MapNotPaired))
	}
	res := make([]Vector, l/2)
	for in, out := 0, 0; in < l; in, out = in+2, out+1 {
		res[out] = Vector{v[in], v[in+1]}
	}
	return Associative(res)
}

// Count returns the number of pairs in the Associative
func (a Associative) Count() int {
	return len(a)
}

// Get returns the Value corresponding to the key in the Associative
func (a Associative) Get(key Value) (Value, bool) {
	l := len(a)
	for i := 0; i < l; i++ {
		mp := a[i]
		k, _ := mp.ElementAt(0)
		if k == key {
			v, _ := mp.ElementAt(1)
			return v, true
		}
	}
	return Null, false
}

// First returns the first pair of the Associative
func (a Associative) First() Value {
	if len(a) > 0 {
		return a[0]
	}
	return Null
}

// Rest returns the pairs of the List that follow the first
func (a Associative) Rest() Sequence {
	if len(a) > 1 {
		return a[1:]
	}
	return EmptyAssociative
}

// IsEmpty returns whether or not this sequence is empty
func (a Associative) IsEmpty() bool {
	return len(a) == 0
}

// Split breaks the Associative into its components (first, rest, ok)
func (a Associative) Split() (Value, Sequence, bool) {
	if len(a) > 0 {
		return a[0], a[1:], true
	}
	return Null, Null, false
}

// Car returns the first element of a Pair
func (a Associative) Car() Value {
	return SequenceCar(a)
}

// Cdr returns the second element of a Pair
func (a Associative) Cdr() Value {
	return SequenceCdr(a)
}

// Prepend inserts a pair at the beginning of the Associative
func (a Associative) Prepend(v Value) Sequence {
	if mp, ok := v.(Vector); ok && mp.Count() == 2 {
		return append(Associative{mp}, a...)
	}
	panic(fmt.Errorf(ExpectedPair))
}

// Caller turns Associative into a callable type
func (a Associative) Caller() Call {
	return makeMappedCall(a)
}

// Convention returns the function's calling convention
func (a Associative) Convention() Convention {
	return ApplicativeCall
}

// CheckArity performs a compile-time arity check for the function
func (a Associative) CheckArity(argCount int) error {
	return checkRangedArity(1, 2, argCount)
}

// String converts this Associative to a string
func (a Associative) String() string {
	var b bytes.Buffer
	l := len(a)

	b.WriteString("{")
	for i := 0; i < l; i++ {
		if i > 0 {
			b.WriteString(", ")
		}
		mp := a[i]
		k, _ := mp.ElementAt(0)
		v, _ := mp.ElementAt(1)
		b.WriteString(MaybeQuoteString(k))
		b.WriteString(" ")
		b.WriteString(MaybeQuoteString(v))
	}
	b.WriteString("}")
	return b.String()
}
