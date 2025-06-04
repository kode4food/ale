package data_test

import (
	"testing"

	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/internal/types"
	"github.com/kode4food/ale/pkg/data"
)

func TestSimpleList(t *testing.T) {
	as := assert.New(t)
	n := F(12)
	l := L(n)
	as.Equal(n, l.Car())
	as.Equal(data.Null, l.Cdr())
}

func TestList(t *testing.T) {
	as := assert.New(t)
	n1 := I(12)
	l1 := data.NewList(n1)

	as.Equal(n1, l1.Car())
	as.Equal(data.Null, l1.Cdr())
	as.True(l1.Cdr().(data.Sequence).IsEmpty())

	n2 := F(20.5)
	l2 := l1.Prepend(n2).(*data.List)

	as.String("()", data.Null)
	as.String("(20.5 12)", l2)
	as.Equal(n2, l2.Car())
	as.Identical(l1, l2.Cdr())
	as.Number(2, l2.Count())

	r, ok := l2.ElementAt(1)
	as.True(ok)
	as.Equal(I(12), r)
	as.Number(2, l2.Count())

	r, ok = data.Null.ElementAt(1)
	as.False(ok)
	as.Equal(data.Null, r)
}

func TestListReverse(t *testing.T) {
	as := assert.New(t)

	l1 := data.NewList(I(1), I(2), I(3), I(4))
	l2 := l1.Reverse().(*data.List)

	as.String("(1 2 3 4)", l1)
	as.Number(4, l1.Count())
	as.String("(4 3 2 1)", l2)
	as.Number(4, l2.Count())

	as.String(`(2 1)`, data.NewList(I(1), I(2)).Reverse())
	as.String("()", data.Null.Reverse())
}

func TestListCaller(t *testing.T) {
	as := assert.New(t)

	l1 := data.NewList(I(99), I(37))
	as.Number(99, l1.Call(I(0)))
	as.Number(37, l1.Call(I(1)))
	as.Nil(l1.Call(I(2)))
	as.String("defaulted", l1.Call(I(2), S("defaulted")))
}

func TestListEquality(t *testing.T) {
	as := assert.New(t)

	l1 := data.NewList(I(99), I(37), I(56))
	l2 := data.NewList(I(99), I(37), I(56))
	l3 := data.NewList(I(99), I(37), I(55))
	l4 := data.NewList()
	l5 := data.NewList()

	as.True(l1.Equal(l1))
	as.True(l1.Equal(l2))
	as.False(l1.Equal(l3))
	as.False(l1.Equal(I(55)))
	as.True(l4.Equal(l5))
	as.False(l4.Equal(I(56)))
}

func TestListAsKey(t *testing.T) {
	as := assert.New(t)

	o1, err := data.ValuesToObject(
		L(S("hello"), S("there")), I(42),
		L(S("hello")), I(96),
		L(S("there")), I(128),
	)

	if as.NoError(err) {
		v, ok := o1.Get(L(S("hello")))
		as.True(ok)
		as.Equal(I(96), v)

		v, ok = o1.Get(L(S("hello"), S("there")))
		as.True(ok)
		as.Equal(I(42), v)
	}
}

func TestEmptyList(t *testing.T) {
	as := assert.New(t)

	l := data.NewList()
	as.Equal(data.Null, l)

	as.Nil(l.Car())
	as.Nil(l.Cdr())

	f, r, ok := l.Split()
	as.Nil(f)
	as.Nil(r)
	as.False(ok)

	as.True(types.BasicNull.Equal(l.Type()))
}

func TestListCall(t *testing.T) {
	as := assert.New(t)

	l1 := data.NewList(I(1), I(2), I(3), I(99))
	l2 := data.Null

	as.Number(2, l1.Call(I(1)))
	as.Nil(l2.Call(I(1)))

	as.String("hello", l1.Call(I(99), S("hello")))
	as.String("hello", l2.Call(I(99), S("hello")))

	testSequenceCallInterface(as, l1)
	testSequenceCallInterface(as, l2)
}

func TestListHashCode(t *testing.T) {
	as := assert.New(t)
	l1 := L(I(8), I(4), I(2), I(1))
	l2 := L(I(8), I(4), I(2), I(1))
	l3 := L(I(4), I(2), I(1))
	as.Equal(l1.HashCode(), l2.HashCode())
	as.NotEqual(l1.HashCode(), l3.HashCode())

	l4 := l3.Prepend(I(8)).(*data.List)
	as.Equal(l1.HashCode(), l4.HashCode())
	as.Equal(l2.HashCode(), l4.HashCode())
}
