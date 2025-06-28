package data

import (
	"math/rand/v2"
	"slices"
	"strings"

	"github.com/kode4food/ale/internal/basics"
	"github.com/kode4food/ale/internal/lang"
	"github.com/kode4food/ale/internal/types"
)

// Vector is a fixed-length array of Values
type Vector []Value

var (
	// EmptyVector represents an empty Vector
	EmptyVector = Vector{}

	vectorSalt = rand.Uint64()

	// compile-time checks for interface implementation
	_ interface {
		Appender
		Hashed
		Prepender
		Procedure
		RandomAccess
		Reverser
		Typed
	} = EmptyVector
)

// NewVector creates a new Vector instance
func NewVector(vals ...Value) Vector {
	return slices.Clone(vals)
}

func (v Vector) Count() int {
	return len(v)
}

func (v Vector) ElementAt(index int) (Value, bool) {
	if index >= 0 && index < len(v) {
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
	res := make(Vector, len(v)+1)
	res[0] = e
	copy(res[1:], v)
	return res
}

func (v Vector) Append(e Value) Sequence {
	res := make(Vector, len(v)+1)
	copy(res, v)
	res[len(v)] = e
	return res
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

func (v Vector) CheckArity(argc int) error {
	return CheckRangedArity(1, 2, argc)
}

func (v Vector) Call(args ...Value) Value {
	res, err := sliceRangedCall(v, args)
	if err != nil {
		panic(err)
	}
	return Vector(res)
}

func (v Vector) Equal(other Value) bool {
	if o, ok := other.(Vector); ok {
		return basics.EqualFunc(v, o, Equal)
	}
	return false
}

func (v Vector) String() string {
	var b strings.Builder
	b.WriteString(lang.VectorStart)
	for i, e := range v {
		if i > 0 {
			b.WriteString(lang.Space)
		}
		b.WriteString(ToQuotedString(e))
	}
	b.WriteString(lang.VectorEnd)
	return b.String()
}

func (Vector) Type() types.Type {
	return types.BasicVector
}

func (v Vector) HashCode() uint64 {
	res := vectorSalt
	for i, e := range v {
		res ^= HashCode(e)
		res ^= HashInt(i)
	}
	return res
}
