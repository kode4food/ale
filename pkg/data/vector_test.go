package data_test

import (
	"testing"

	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/pkg/data"
)

func TestVector(t *testing.T) {
	as := assert.New(t)

	v1 := data.NewVector(S("hello"), S("how"), S("are"), S("you?"))
	as.Number(4, v1.Count())
	as.String("hello", v1.Car())
	as.Number(3, v1.Cdr().(data.Counted).Count())

	r, ok := v1.ElementAt(2)
	as.True(ok)
	as.String("are", r)
	as.String(`["hello" "how" "are" "you?"]`, v1)

	v2 := v1.Prepend(S("oh")).(data.Vector)
	as.Number(5, v2.Count())
	as.Number(4, v1.Count())

	v3 := v2.Append(S("good?")).(data.Vector)
	r, ok = v3.ElementAt(5)
	as.True(ok)
	as.String("good?", r)
	as.Number(6, v3.Count())

	r, ok = v3.ElementAt(0)
	as.True(ok)
	as.String("oh", r)

	r, ok = v3.ElementAt(3)
	as.True(ok)
	as.String("are", r)
}

func TestVectorReverse(t *testing.T) {
	as := assert.New(t)

	as.String("[4 3 2 1]", data.NewVector(I(1), I(2), I(3), I(4)).Reverse())
	as.String("[]", data.EmptyVector.Reverse())
}

func TestEmptyVector(t *testing.T) {
	as := assert.New(t)

	v := data.EmptyVector
	as.Nil(v.Car())
	as.String("[]", v)
	as.String("[]", v.Cdr())
}

func TestVectorCall(t *testing.T) {
	as := assert.New(t)

	v1 := data.NewVector(S("hello"), S("how"), S("are"), S("you?"))
	as.Equal(v1, v1.Call(I(0)))
	as.Equal(v1[1:], v1.Call(I(1)))
	as.Equal(data.EmptyVector, v1.Call(I(4)))

	testSequenceCallInterface(as, v1)
}

func TestVectorEquality(t *testing.T) {
	as := assert.New(t)

	v1 := data.NewVector(S("hello"), S("how"), S("are"), S("you?"))
	v2 := data.NewVector(S("hello"), S("how"), S("are"), S("you?"))
	v3 := data.NewVector(S("hello"), S("are"), S("you?"), S("how"))
	v4 := data.NewVector(S("hello"), S("how"), S("are"))

	as.True(v1.Equal(v1))
	as.True(v1.Equal(v2))
	as.False(v1.Equal(v3))
	as.False(v1.Equal(v4))
	as.False(v1.Equal(I(32)))
}

func TestVectorSplit(t *testing.T) {
	as := assert.New(t)

	v1 := data.NewVector(S("hello"), S("how"), S("are"), S("you?"))
	f, r, ok := v1.Split()

	as.True(ok)
	as.Equal(S("hello"), f)
	as.Equal(
		data.Vector{S("how"), S("are"), S("you?")},
		r.(data.Vector),
	)

	v2 := data.NewVector(S("hello"))
	f, r, ok = v2.Split()
	as.True(ok)
	as.Equal(S("hello"), f)
	as.Equal(data.Vector{}, r.(data.Vector))

	v3 := data.NewVector()
	_, _, ok = v3.Split()
	as.False(ok)
}

func TestVectorAsKey(t *testing.T) {
	as := assert.New(t)

	o1, err := data.ValuesToObject(
		V(S("hello"), S("there")), I(42),
		V(S("hello")), I(96),
		V(S("there")), I(128),
	)

	if as.NoError(err) {
		v, ok := o1.Get(V(S("hello")))
		as.True(ok)
		as.Equal(I(96), v)

		v, ok = o1.Get(V(S("hello"), S("there")))
		as.True(ok)
		as.Equal(I(42), v)
	}
}

func TestVectorAppendIsolation(t *testing.T) {
	as := assert.New(t)

	v1 := data.NewVector(I(1), I(2), I(3))
	v2 := v1.Append(I(4)).(data.Vector).Append(I(5))
	v3 := v1.Append(I(6)).(data.Vector).Append(I(7))

	as.String("[1 2 3]", v1)
	as.String("[1 2 3 4 5]", v2)
	as.String("[1 2 3 6 7]", v3)
}
