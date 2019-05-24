package stdlib_test

import (
	"testing"

	"gitlab.com/kode4food/ale/data"
	"gitlab.com/kode4food/ale/internal/assert"
	. "gitlab.com/kode4food/ale/internal/assert/helpers"
	"gitlab.com/kode4food/ale/stdlib"
)

func TestFilter(t *testing.T) {
	as := assert.New(t)

	filterTest := func(args ...data.Value) data.Value {
		return B(string(args[0].(data.String)) != "filtered out")
	}

	l := L(S("first"), S("filtered out"), S("last"))
	w := stdlib.Filter(l, filterTest)

	v1 := w.First()
	as.String("first", v1)

	v2 := w.Rest().First()
	as.String("last", v2)

	r1 := w.Rest().Rest()
	as.True(r1.IsEmpty())

	p := w.Prepend(S("filtered out"))
	v4 := p.First()
	r2 := p.Rest()
	as.String("filtered out", v4)
	as.Equal(w.First(), r2.First())
}

func TestFiltered(t *testing.T) {
	as := assert.New(t)

	l := L(S("first"), S("middle"), S("last"))
	fn1 := func(args ...data.Value) data.Value {
		return B(string(args[0].(data.String)) != "middle")
	}
	w1 := stdlib.Filter(l, fn1)
	v1 := w1.First()
	as.String("first", v1)

	v2 := w1.Rest().First()
	as.String("last", v2)

	r1 := w1.Rest().Rest()
	as.True(r1.IsEmpty())
}

func TestConcat(t *testing.T) {
	as := assert.New(t)

	l1 := L(S("first"), S("middle"), S("last"))
	l2 := data.EmptyList
	l3 := V(I(1), I(2), I(3))
	l4 := L(S("blah1"), S("blah2"), S("blah3"))
	l5 := data.EmptyList

	w1 := stdlib.Concat(l1, l2, l3, l4, l5)
	expect := `("first" "middle" "last" 1 2 3 "blah1" "blah2" "blah3")`
	as.String(expect, data.MakeSequenceStr(w1))
}
