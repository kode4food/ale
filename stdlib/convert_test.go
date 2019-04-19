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
	a1 := stdlib.SequenceToAssociative(l3)
	a2 := stdlib.SequenceToAssociative(a1)

	l4 := L(S("hello"), data.Nil, S("there"), v1)
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

func identity(args ...data.Value) data.Value {
	return args[0]
}

func TestUncountedConversions(t *testing.T) {
	as := assert.New(t)
	l1 := stdlib.Map(L(S("hello"), S("there")), identity)
	v1 := stdlib.SequenceToVector(l1)
	v2 := stdlib.SequenceToVector(v1)
	l2 := stdlib.SequenceToList(stdlib.Map(v2, identity))
	l3 := stdlib.SequenceToList(l2)
	a1 := stdlib.SequenceToAssociative(stdlib.Map(l3, identity))
	a2 := stdlib.SequenceToAssociative(a1)

	l4 := stdlib.Map(L(S("hello"), data.Nil, S("there"), v1), identity)
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
	defer as.ExpectPanic(data.ExpectedPair)
	stdlib.SequenceToAssociative(v1)
}

func TestUncountedAssocConvertError(t *testing.T) {
	as := assert.New(t)

	v1 := stdlib.Map(V(K("boom")), identity)
	defer as.ExpectPanic(data.ExpectedPair)
	stdlib.SequenceToAssociative(v1)
}
