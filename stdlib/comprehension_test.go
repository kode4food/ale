package stdlib_test

import (
	"testing"

	"gitlab.com/kode4food/ale/data"
	"gitlab.com/kode4food/ale/internal/assert"
	. "gitlab.com/kode4food/ale/internal/assert/helpers"
	"gitlab.com/kode4food/ale/stdlib"
)

func TestMap(t *testing.T) {
	as := assert.New(t)

	concatTest := func(args ...data.Value) data.Value {
		return S("this is the " + string(args[0].(data.String)))
	}

	l := L(S("first"), S("middle"), S("last"))
	w := stdlib.Map(l, concatTest)

	v1 := w.First()
	as.String("this is the first", v1)

	v2 := w.Rest().First()
	as.String("this is the middle", v2)

	v3 := w.Rest().Rest().First()
	as.String("this is the last", v3)

	r1 := w.Rest().Rest().Rest()
	as.True(r1.IsEmpty())

	p1 := w.Prepend(S("not mapped"))
	p2 := p1.Prepend(S("also not mapped"))
	v4 := p1.First()
	r2 := p1.Rest()

	as.String("not mapped", v4)
	as.Equal(w.First(), r2.First())
	as.String("also not mapped", p2.First())
}

func TestMapParallel(t *testing.T) {
	as := assert.New(t)

	addTest := func(args ...data.Value) data.Value {
		return args[0].(data.Integer) + args[1].(data.Integer)
	}

	s1 := L(I(1), I(2), I(3), I(4))
	s2 := V(I(5), I(10), I(15), I(20), I(30))

	w := stdlib.MapParallel(V(s1, s2), addTest)

	as.Float(6, w.First())
	as.Float(12, w.Rest().First())
	as.Float(18, w.Rest().Rest().First())
	as.Float(24, w.Rest().Rest().Rest().First())

	s3 := w.Rest().Rest().Rest().Rest()
	as.True(s3.IsEmpty())
}

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

func TestFilteredAndMapped(t *testing.T) {
	as := assert.New(t)

	l := L(S("first"), S("middle"), S("last"))
	fn1 := func(args ...data.Value) data.Value {
		return B(string(args[0].(data.String)) != "middle")
	}
	w1 := stdlib.Filter(l, fn1)

	fn2 := func(args ...data.Value) data.Value {
		return S("this is the " + string(args[0].(data.String)))
	}
	w2 := stdlib.Map(w1, fn2)

	v1 := w2.First()
	as.String("this is the first", v1)

	v2 := w2.Rest().First()
	as.String("this is the last", v2)

	r1 := w2.Rest().Rest()
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

func TestReduce(t *testing.T) {
	as := assert.New(t)

	add := func(args ...data.Value) data.Value {
		return args[0].(data.Integer) + args[1].(data.Integer)
	}

	as.Integer(30, stdlib.Reduce(V(I(10), I(20)), add))
	as.Integer(60, stdlib.Reduce(V(I(10), I(20), I(30)), add))
	as.Integer(100, stdlib.Reduce(V(I(10), I(20), I(30), I(40)), add))
}

func TestTakeDrop(t *testing.T) {
	as := assert.New(t)

	s1 := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}
	v1 := data.EmptyVector
	for _, e := range s1 {
		v1 = v1.Append(S(e)).(data.Vector)
	}

	t1 := stdlib.Take(v1, 4)
	t2 := t1.Prepend(S("0"))
	d1 := stdlib.Drop(v1, 4)
	d2 := d1.Prepend(S("4"))
	t3 := stdlib.Take(d1, 6)
	d3 := stdlib.Drop(t3, 6)
	d4 := stdlib.Drop(t3, 8)

	as.String(`("1" "2" "3" "4")`, data.MakeSequenceStr(t1))
	as.String(`("0" "1" "2" "3" "4")`, data.MakeSequenceStr(t2))
	as.String(`("5" "6" "7" "8" "9" "10")`, data.MakeSequenceStr(d1))
	as.String(`("4" "5" "6" "7" "8" "9" "10")`, data.MakeSequenceStr(d2))
	as.String(`("5" "6" "7" "8" "9" "10")`, data.MakeSequenceStr(t3))
	as.True(d3.IsEmpty())
	as.True(d4.IsEmpty())
}
