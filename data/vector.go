package data

import (
	"fmt"
	"math/rand/v2"
	"slices"
	"strings"

	"github.com/kode4food/ale"
	"github.com/kode4food/ale/internal/basics"
	"github.com/kode4food/ale/internal/lang"
	"github.com/kode4food/ale/internal/types"
)

// Vector is a fixed-length array of Values
type Vector []ale.Value

var (
	// EmptyVector represents an empty Vector
	EmptyVector = Vector{}

	vectorSalt = rand.Uint64()

	// compile-time checks for interface implementation
	_ interface {
		Appender
		Hashed
		Indexed
		Prepender
		Procedure
		Reverser
		ale.Typed
		fmt.Stringer
	} = EmptyVector
)

// NewVector creates a new Vector instance
func NewVector(vals ...ale.Value) Vector {
	return slices.Clone(vals)
}

func (v Vector) Count() int {
	return len(v)
}

func (v Vector) ElementAt(index int) (ale.Value, bool) {
	if index >= 0 && index < len(v) {
		return v[index], true
	}
	return Null, false
}

func (v Vector) IsEmpty() bool {
	return len(v) == 0
}

func (v Vector) Car() ale.Value {
	if len(v) > 0 {
		return v[0]
	}
	return Null
}

func (v Vector) Cdr() ale.Value {
	if len(v) > 1 {
		return v[1:]
	}
	return EmptyVector
}

func (v Vector) Split() (ale.Value, Sequence, bool) {
	switch len(v) {
	case 0:
		return Null, EmptyVector, false
	case 1:
		return v[0], EmptyVector, true
	default:
		return v[0], v[1:], true
	}
}

func (v Vector) Prepend(e ale.Value) Sequence {
	res := make(Vector, len(v)+1)
	res[0] = e
	copy(res[1:], v)
	return res
}

func (v Vector) Append(e ale.Value) Sequence {
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

func (v Vector) IndexOf(val ale.Value) (int, bool) {
	i := slices.IndexFunc(v, val.Equal)
	return i, i != -1
}

func (v Vector) CheckArity(argc int) error {
	return CheckRangedArity(1, 2, argc)
}

func (v Vector) Call(args ...ale.Value) ale.Value {
	res, err := sliceRangedCall(v, args)
	if err != nil {
		panic(err)
	}
	return Vector(res)
}

func (v Vector) Equal(other ale.Value) bool {
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

func (v Vector) Type() ale.Type {
	return types.MakeLiteral(types.BasicVector, v)
}

func (v Vector) HashCode() uint64 {
	res := vectorSalt
	for i, e := range v {
		res ^= HashCode(e)
		res ^= HashInt(i)
	}
	return res
}
