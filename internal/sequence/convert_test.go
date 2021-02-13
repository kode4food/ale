package sequence_test

import (
	"testing"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/internal/sequence"
)

func TestSequenceConversions(t *testing.T) {
	as := assert.New(t)
	l1 := L(S("hello"), S("there"))
	v1 := sequence.ToVector(l1)
	v2 := sequence.ToVector(v1)
	l2 := sequence.ToList(v2)
	l3 := sequence.ToList(l2)

	a1, err := sequence.ToObject(l3)
	as.NotNil(a1)
	as.Nil(err)

	a2, err := sequence.ToObject(a1)
	as.NotNil(a2)
	as.Nil(err)

	l4 := L(S("hello"), data.Nil, S("there"), v1)
	s1 := sequence.ToStr(l4)
	s2 := sequence.ToStr(s1)

	as.String(`["hello" "there"]`, v1)
	as.Identical(v1, v2)
	as.String(`("hello" "there")`, l2)
	as.Identical(l2, l3)
	as.String(`{"hello" "there"}`, a1)
	as.Identical(a1, a2)
	as.String(`hellothere["hello" "there"]`, s1)
	as.Identical(s1, s2)
}

var alwaysTrue = data.Applicative(func(_ ...data.Value) data.Value {
	return data.True
}, 1)

func TestUncountedConversions(t *testing.T) {
	as := assert.New(t)
	l1 := sequence.Filter(L(S("hello"), S("there")), alwaysTrue)
	v1 := sequence.ToVector(l1)
	v2 := sequence.ToVector(v1)
	l2 := sequence.ToList(sequence.Filter(v2, alwaysTrue))
	l3 := sequence.ToList(l2)

	a1, err := sequence.ToObject(sequence.Filter(l3, alwaysTrue))
	as.NotNil(a1)
	as.Nil(err)

	a2, err := sequence.ToObject(a1)
	as.NotNil(a1)
	as.Nil(err)

	l4 := sequence.Filter(L(S("hello"), data.Nil, S("there"), v1), alwaysTrue)
	s1 := sequence.ToStr(l4)

	as.String(`["hello" "there"]`, v1)
	as.Identical(v1, v2)
	as.String(`("hello" "there")`, l2)
	as.Identical(l2, l3)
	as.String(`{"hello" "there"}`, a1)
	as.Identical(a1, a2)
	as.String(`hellothere["hello" "there"]`, s1)
}

func TestMappedSequenceError(t *testing.T) {
	as := assert.New(t)

	v1 := V(K("boom"))
	o, err := sequence.ToObject(v1)
	as.Nil(o)
	as.Error(err, data.ErrMapNotPaired)
}

func TestUncountedMappedSequenceError(t *testing.T) {
	as := assert.New(t)

	v1 := sequence.Filter(V(K("boom")), alwaysTrue)
	o, err := sequence.ToObject(v1)
	as.Nil(o)
	as.Error(err, data.ErrMapNotPaired)
}
