package data

import (
	"bytes"
	"math/rand"

	"github.com/kode4food/ale/types"
)

// Vector is a fixed-length array of Values
type Vector Values

// EmptyVector represents an empty Vector
var (
	EmptyVector = Vector{}

	vectorHash = rand.Uint64()
)

// NewVector creates a new Vector instance
func NewVector(v ...Value) Vector {
	return Vector(v)
}

func (v Vector) Values() Values {
	return Values(v)
}

func (v Vector) Count() int {
	return len(v)
}

func (v Vector) ElementAt(index int) (Value, bool) {
	if index >= 0 && index < len(v) {
		return v[index], true
	}
	return Nil, false
}

func (v Vector) IsEmpty() bool {
	return len(v) == 0
}

func (v Vector) Car() Value {
	if len(v) > 0 {
		return v[0]
	}
	return Nil
}

func (v Vector) Cdr() Value {
	if len(v) > 1 {
		return v[1:]
	}
	return EmptyVector
}

func (v Vector) Split() (Value, Sequence, bool) {
	lv := len(v)
	if lv > 1 {
		return v[0], v[1:], true
	} else if lv == 1 {
		return v[0], EmptyVector, true
	}
	return Nil, EmptyVector, false
}

func (v Vector) Prepend(e Value) Sequence {
	return append(Vector{e}, v...)
}

func (v Vector) Append(e Value) Sequence {
	return append(v, e)
}

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

func (v Vector) Call(args ...Value) Value {
	return indexedCall(v, args)
}

func (v Vector) Convention() Convention {
	return ApplicativeCall
}

func (v Vector) CheckArity(argCount int) error {
	return checkRangedArity(1, 2, argCount)
}

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

func (v Vector) String() string {
	var b bytes.Buffer
	b.WriteString("[")
	for i, e := range v {
		if i > 0 {
			b.WriteString(" ")
		}
		b.WriteString(MaybeQuoteString(e))
	}
	b.WriteString("]")
	return b.String()
}

func (Vector) Type() types.Type {
	return types.AnyVector
}

func (v Vector) HashCode() uint64 {
	h := vectorHash
	for _, e := range v {
		h *= HashCode(e)
	}
	return h
}
