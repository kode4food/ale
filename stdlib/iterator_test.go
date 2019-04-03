package stdlib_test

import (
	"testing"

	"gitlab.com/kode4food/ale/api"
	"gitlab.com/kode4food/ale/internal/assert"
	. "gitlab.com/kode4food/ale/internal/assert/helpers"
	"gitlab.com/kode4food/ale/stdlib"
)

func TestListIterate(t *testing.T) {
	as := assert.New(t)
	n1 := F(12)
	l1 := L(n1)
	as.Equal(n1, l1.First())
	as.Identical(api.EmptyList, l1.Rest())
	as.False(l1.Rest().IsSequence())

	n2 := F(20.5)
	l2 := l1.Conjoin(n2)
	as.Equal(n2, l2.First())
	as.Identical(l1, l2.Rest())

	sum := F(0.0)
	i := stdlib.Iterate(l2)
	for {
		v, ok := i.Next()
		if !ok {
			break
		}
		sum = sum + v.(api.Float)
	}

	as.Float(32.5, sum)
}

func TestVectorIterate(t *testing.T) {
	as := assert.New(t)

	v := V(S("hello"), S("how"), S("are"), S("you?"))
	i := stdlib.Iterate(v)
	e1, _ := i.Next()
	s1 := i.Rest()
	e2, _ := i.Next()
	s2 := i.Rest()
	e3, _ := i.Next()
	e4, _ := i.Next()
	e5, ok := i.Next()

	as.String("hello", e1)
	as.String("how", e2)
	as.String("are", e3)
	as.String("you?", e4)

	as.Integer(3, api.Count(s1))
	as.Integer(2, api.Count(s2))

	as.Equal(api.Nil, e5)
	as.False(ok)
}

func getTestMap() api.Associative {
	return A(
		V(K("name"), S("Ale")),
		V(K("age"), F(99)),
		V(S("string"), S("value")),
	)
}

func TestAssociativeIterate(t *testing.T) {
	as := assert.New(t)
	m1 := getTestMap()

	i := stdlib.Iterate(m1)
	if v, ok := i.Next(); ok {
		vec := v.(api.Vector)
		key, _ := vec.ElementAt(0)
		val, _ := vec.ElementAt(1)
		as.Equal(K("name"), key)
		as.String("Ale", val)
	} else {
		as.Fail("couldn't get first element")
	}

	if v, ok := i.Next(); ok {
		vec := v.(api.Vector)
		key, _ := vec.ElementAt(0)
		val, _ := vec.ElementAt(1)
		as.Equal(K("age"), key)
		as.Float(99, val)
	} else {
		as.Fail("couldn't get second element")
	}
}
