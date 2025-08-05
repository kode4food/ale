package data_test

import (
	"testing"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
)

func TestTypePredicateCall(t *testing.T) {
	as := assert.New(t)
	l1 := L(I(1), I(2), I(3))
	pred := data.TypePredicateOf(l1)
	if as.NotNil(pred) {
		l2 := L(I(9))
		v1 := V(I(10))
		as.True(pred.Call(l1))
		as.True(pred.Call(L(I(1), I(2), I(3))))
		as.False(pred.Call(l2))
		as.False(pred.Call(v1))
	}
}

func TestPredicateOf(t *testing.T) {
	as := assert.New(t)

	l := L(I(0))
	o := O()

	list := data.TypePredicateOf(l)
	obj := data.TypePredicateOf(o)
	union := data.TypePredicateOf(l, o)

	as.String("list((0))", list)
	as.String("object({})", obj)
	as.String("union(list((0)),object({}))", union)

	as.True(list.Call(l))
	as.False(obj.Call(l))
	as.True(obj.Call(o))
	as.False(obj.Call(l))
	as.True(union.Call(l))
	as.True(union.Call(o))
	as.False(union.Call(data.True))

	as.String("type-predicate(list((0)))", list.Type().Name())
}

func TestTypePredicateEqual(t *testing.T) {
	as := assert.New(t)
	l1 := L(I(1), I(2), I(3))
	l2 := L(I(9))
	l3 := L(I(1), I(2), I(3))
	v1 := V(I(10))
	p1 := data.TypePredicateOf(l1)
	p2 := data.TypePredicateOf(l2)
	p3 := data.TypePredicateOf(v1)
	as.True(l1.Equal(l1))
	as.True(l1.Equal(l3))
	as.False(p1.Equal(p2))
	as.False(p1.Equal(p3))
	as.True(p1.Equal(p1))
}
