package data_test

import (
	"testing"

	"gitlab.com/kode4food/ale/data"
	"gitlab.com/kode4food/ale/internal/assert"
	. "gitlab.com/kode4food/ale/internal/assert/helpers"
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
	l1 := L(n1)

	as.Equal(n1, l1.First())
	as.Equal(data.EmptyList, l1.Rest())
	as.True(l1.Rest().IsEmpty())

	n2 := F(20.5)
	l2 := l1.Prepend(n2).(*data.List)

	as.String("()", data.EmptyList)
	as.String("(20.5 12)", l2)
	as.Equal(n2, l2.First())
	as.Identical(l1, l2.Rest())
	as.Number(2, l2.Count())

	r, ok := l2.ElementAt(1)
	as.True(ok)
	as.Equal(I(12), r)
	as.Number(2, data.Count(l2))

	r, ok = data.EmptyList.ElementAt(1)
	as.False(ok)
	as.Equal(data.Nil, r)
}

func TestListReverse(t *testing.T) {
	as := assert.New(t)

	as.String("(4 3 2 1)", data.NewList(I(1), I(2), I(3), I(4)).Reverse())
	as.String("()", data.EmptyList.Reverse())
}

func TestListCaller(t *testing.T) {
	as := assert.New(t)

	l1 := L(I(99), I(37))
	c1 := l1.Caller()
	as.Number(99, c1(I(0)))
	as.Number(37, c1(I(1)))
	as.Nil(c1(I(2)))
	as.String("defaulted", c1(I(2), S("defaulted")))
}
