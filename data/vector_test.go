package data_test

import (
	"testing"

	"gitlab.com/kode4food/ale/data"
	"gitlab.com/kode4food/ale/internal/assert"
	. "gitlab.com/kode4food/ale/internal/assert/helpers"
)

func TestVector(t *testing.T) {
	as := assert.New(t)

	v1 := V(S("hello"), S("how"), S("are"), S("you?"))
	as.Integer(4, v1.Count())
	as.Integer(4, data.Count(v1))
	as.String("hello", v1.First())
	as.Integer(3, data.Count(v1.Rest()))

	r, ok := v1.ElementAt(2)
	as.True(ok)
	as.String("are", r)
	as.String(`["hello" "how" "are" "you?"]`, v1)

	v2 := v1.Prepend(S("oh")).(data.Vector)
	as.Integer(5, v2.Count())
	as.Integer(4, v1.Count())

	v3 := v2.Append(S("good?")).(data.Vector)
	r, ok = v3.ElementAt(5)
	as.True(ok)
	as.String("good?", r)
	as.Integer(6, v3.Count())

	r, ok = v3.ElementAt(0)
	as.True(ok)
	as.String("oh", r)

	r, ok = v3.ElementAt(3)
	as.True(ok)
	as.String("are", r)
}

func TestEmptyVector(t *testing.T) {
	as := assert.New(t)

	v := data.EmptyVector
	as.Nil(v.First())
	as.String("[]", v)
	as.String("[]", v.Rest())
}

func TestVectorCaller(t *testing.T) {
	as := assert.New(t)

	v1 := V(S("hello"), S("how"), S("are"), S("you?"))
	c1 := v1.Caller()
	as.String("hello", c1(I(0)))
	as.String("how", c1(I(1)))
	as.Nil(c1(I(4)))
	as.String("defaulted", c1(I(4), S("defaulted")))
}
