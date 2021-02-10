package data_test

import (
	"testing"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
)

func TestVector(t *testing.T) {
	as := assert.New(t)

	v1 := V(S("hello"), S("how"), S("are"), S("you?"))
	as.Number(4, v1.Count())
	as.String("hello", v1.First())
	as.Number(3, v1.Rest().(data.Counted).Count())

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

	as.String("[4 3 2 1]", V(I(1), I(2), I(3), I(4)).Reverse())
	as.String("[]", data.EmptyVector.Reverse())
}

func TestEmptyVector(t *testing.T) {
	as := assert.New(t)

	v := data.EmptyVector
	as.Nil(v.First())
	as.String("[]", v)
	as.String("[]", v.Rest())
}

func TestVectorAsFunction(t *testing.T) {
	as := assert.New(t)

	v1 := V(S("hello"), S("how"), S("are"), S("you?")).(data.Function)
	as.String("hello", v1.Call(I(0)))
	as.String("how", v1.Call(I(1)))
	as.Nil(v1.Call(I(4)))
	as.String("defaulted", v1.Call(I(4), S("defaulted")))
}

func TestVectorEquality(t *testing.T) {
	as := assert.New(t)

	v1 := V(S("hello"), S("how"), S("are"), S("you?"))
	v2 := V(S("hello"), S("how"), S("are"), S("you?"))
	v3 := V(S("hello"), S("are"), S("you?"), S("how"))
	v4 := V(S("hello"), S("how"), S("are"))

	as.True(v1.Equal(v1))
	as.True(v1.Equal(v2))
	as.False(v1.Equal(v3))
	as.False(v1.Equal(v4))
	as.False(v1.Equal(I(32)))
}
