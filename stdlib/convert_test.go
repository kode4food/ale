package stdlib_test

import (
	"testing"

	"gitlab.com/kode4food/ale/data"
	"gitlab.com/kode4food/ale/internal/assert"
	. "gitlab.com/kode4food/ale/internal/assert/helpers"
	"gitlab.com/kode4food/ale/stdlib"
)

func TestSequenceConversions(t *testing.T) {
	as := assert.New(t)
	l1 := L(S("hello"), S("there"))
	v1 := stdlib.SequenceToVector(l1)
	v2 := stdlib.SequenceToVector(v1)
	l2 := stdlib.SequenceToList(v2)
	l3 := stdlib.SequenceToList(l2)
	a1 := stdlib.SequenceToObject(l3)
	a2 := stdlib.SequenceToObject(a1)

	l4 := L(S("hello"), data.Null, S("there"), v1)
	s1 := stdlib.SequenceToStr(l4)
	s2 := stdlib.SequenceToStr(s1)

	as.String(`["hello" "there"]`, v1)
	as.Identical(v1, v2)
	as.String(`("hello" "there")`, l2)
	as.Identical(l2, l3)
	as.String(`{"hello" "there"}`, a1)
	as.Identical(a1, a2)
	as.String(`hellothere["hello" "there"]`, s1)
	as.Identical(s1, s2)
}

func alwaysTrue(_ ...data.Value) data.Value {
	return data.True
}

func TestUncountedConversions(t *testing.T) {
	as := assert.New(t)
	l1 := stdlib.Filter(L(S("hello"), S("there")), alwaysTrue)
	v1 := stdlib.SequenceToVector(l1)
	v2 := stdlib.SequenceToVector(v1)
	l2 := stdlib.SequenceToList(stdlib.Filter(v2, alwaysTrue))
	l3 := stdlib.SequenceToList(l2)
	a1 := stdlib.SequenceToObject(stdlib.Filter(l3, alwaysTrue))
	a2 := stdlib.SequenceToObject(a1)

	l4 := stdlib.Filter(L(S("hello"), data.Null, S("there"), v1), alwaysTrue)
	s1 := stdlib.SequenceToStr(l4)

	as.String(`["hello" "there"]`, v1)
	as.Identical(v1, v2)
	as.String(`("hello" "there")`, l2)
	as.Identical(l2, l3)
	as.String(`{"hello" "there"}`, a1)
	as.Identical(a1, a2)
	as.String(`hellothere["hello" "there"]`, s1)
}

func TestAssocConvertError(t *testing.T) {
	as := assert.New(t)

	v1 := V(K("boom"))
	defer as.ExpectPanic(data.ObjectNotPaired)
	stdlib.SequenceToObject(v1)
}

func TestUncountedAssocConvertError(t *testing.T) {
	as := assert.New(t)

	v1 := stdlib.Filter(V(K("boom")), alwaysTrue)
	defer as.ExpectPanic(data.ObjectNotPaired)
	stdlib.SequenceToObject(v1)
}
