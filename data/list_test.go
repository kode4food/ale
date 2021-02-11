package data_test

import (
	"testing"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
)

func TestSimpleList(t *testing.T) {
	as := assert.New(t)
	n := F(12)
	l := L(n)
	as.Equal(n, l.First())
	as.Equal(data.EmptyList, l.Rest())
}

func TestList(t *testing.T) {
	as := assert.New(t)
	n1 := I(12)
	l1 := data.NewList(n1)

	as.Equal(n1, l1.First())
	as.Equal(data.EmptyList, l1.Rest())
	as.True(l1.Rest().IsEmpty())

	n2 := F(20.5)
	l2 := l1.Prepend(n2).(data.List)

	as.String("()", data.EmptyList)
	as.String("(20.5 12)", l2)
	as.Equal(n2, l2.First())
	as.Identical(l1, l2.Rest())
	as.Number(2, l2.Count())

	r, ok := l2.ElementAt(1)
	as.True(ok)
	as.Equal(I(12), r)
	as.Number(2, l2.Count())

	r, ok = data.EmptyList.ElementAt(1)
	as.False(ok)
	as.Equal(data.Nil, r)
}

func TestListReverse(t *testing.T) {
	as := assert.New(t)

	l1 := data.NewList(I(1), I(2), I(3), I(4))
	l2 := l1.Reverse().(data.List)

	as.String("(1 2 3 4)", l1)
	as.Number(4, l1.Count())
	as.String("(4 3 2 1)", l2)
	as.Number(4, l2.Count())

	as.String(`(2 1)`, data.NewList(I(1), I(2)).Reverse())
	as.String("()", data.EmptyList.Reverse())
}

func TestListCaller(t *testing.T) {
	as := assert.New(t)

	l1 := data.NewList(I(99), I(37))
	c1 := l1.(data.Function)
	as.Number(99, c1.Call(I(0)))
	as.Number(37, c1.Call(I(1)))
	as.Nil(c1.Call(I(2)))
	as.String("defaulted", c1.Call(I(2), S("defaulted")))
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
