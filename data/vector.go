package data

import (
	"bytes"
	"math/rand"
	"slices"

	"github.com/kode4food/ale/internal/types"
)

// Vector is a fixed-length array of Values
type Vector []Value

// EmptyVector represents an empty Vector
var (
	EmptyVector = Vector{}

	vectorHash = rand.Uint64()

	// compile-time checks for interface implementation
	_ Appender     = EmptyVector
	_ Caller       = EmptyVector
	_ Hashed       = EmptyVector
	_ Prepender    = EmptyVector
	_ RandomAccess = EmptyVector
	_ Reverser     = EmptyVector
	_ Typed        = EmptyVector
)

// NewVector creates a new Vector instance
func NewVector(v ...Value) Vector {
	return v
}

func (v Vector) Count() Integer {
	return Integer(len(v))
}

func (v Vector) ElementAt(index Integer) (Value, bool) {
	if index >= 0 && index < Integer(len(v)) {
		return v[index], true
	}
	return Null, false
}

func (v Vector) IsEmpty() bool {
	return len(v) == 0
}

func (v Vector) Car() Value {
	if len(v) > 0 {
		return v[0]
	}
	return Null
}

func (v Vector) Cdr() Value {
	if len(v) > 1 {
		return v[1:]
	}
	return EmptyVector
}

func (v Vector) Split() (Value, Sequence, bool) {
	switch len(v) {
	case 0:
		return Null, EmptyVector, false
	case 1:
		return v[0], EmptyVector, true
	default:
		return v[0], v[1:], true
	}
}

func (v Vector) Prepend(e Value) Sequence {
	res := make(Vector, 1, len(v)+1)
	res[0] = e
	return append(res, v...)
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

func (v Vector) IndexOf(val Value) (int, bool) {
	i := slices.IndexFunc(v, val.Equal)
	return i, i != -1
}

func (v Vector) Call(args ...Value) Value {
	return indexedCall(v, args)
}

func (v Vector) CheckArity(argCount int) error {
	return checkRangedArity(1, 2, argCount)
}

func (v Vector) Equal(other Value) bool {
	if o, ok := other.(Vector); ok {
		return slices.EqualFunc(v, o, Equal)
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
		b.WriteString(ToQuotedString(e))
	}
	b.WriteString("]")
	return b.String()
}

func (Vector) Type() types.Type {
	return types.BasicVector
}

func (v Vector) HashCode() uint64 {
	res := vectorHash
	for i, e := range v {
		res ^= HashCode(e)
		res ^= HashInt(i)
	}
	return res
}
