package data_test

import (
	"testing"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
)

func TestTypePredicateCall(t *testing.T) {
	as := assert.New(t)
	l1 := data.NewList(I(1), I(2), I(3))
	pred := data.TypePredicateOf(l1)
	as.NotNil(pred)

	l2 := data.NewList(I(9))
	v1 := data.NewVector(I(10))
	as.True(pred.Call(l1))
	as.True(pred.Call(l2))
	as.False(pred.Call(v1))
}

func TestPredicateOf(t *testing.T) {
	as := assert.New(t)

	l := L(I(0))
	o := O()

	list := data.TypePredicateOf(l)
	obj := data.TypePredicateOf(o)
	union := data.TypePredicateOf(l, o)

	as.Equal(LS("list"), list.Name())
	as.Equal(LS("object"), obj.Name())
	as.Equal(LS("union(list,object)"), union.Name())

	as.True(list.Call(l))
	as.False(obj.Call(l))
	as.True(obj.Call(o))
	as.False(obj.Call(l))
	as.True(union.Call(l))
	as.True(union.Call(o))
	as.False(union.Call(data.True))
}

func TestTypePredicateEqual(t *testing.T) {
	as := assert.New(t)
	l1 := data.NewList(I(1), I(2), I(3))
	l2 := data.NewList(I(9))
	v1 := data.NewVector(I(10))
	p1 := data.TypePredicateOf(l1)
	p2 := data.TypePredicateOf(l2)
	p3 := data.TypePredicateOf(v1)
	as.True(p1.Equal(p2))
	as.False(p1.Equal(p3))
}
