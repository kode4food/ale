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
	as.False(l1.Rest().IsSequence())

	n2 := F(20.5)
	l2 := l1.Prepend(n2).(*data.List)

	as.String("()", data.EmptyList)
	as.String("(20.5 12)", l2)
	as.Equal(n2, l2.First())
	as.Identical(l1, l2.Rest())
	as.Integer(2, l2.Count())

	r, ok := l2.ElementAt(1)
	as.True(ok)
	as.Equal(I(12), r)
	as.Integer(2, data.Count(l2))

	r, ok = data.EmptyList.ElementAt(1)
	as.False(ok)
	as.Equal(data.Nil, r)
}
