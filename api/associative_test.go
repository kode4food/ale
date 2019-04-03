package api_test

import (
	"testing"

	"gitlab.com/kode4food/ale/api"
	"gitlab.com/kode4food/ale/internal/assert"
	. "gitlab.com/kode4food/ale/internal/assert/helpers"
)

func getTestMap() api.Associative {
	return A(
		V(K("name"), S("Ale")),
		V(K("age"), I(99)),
		V(S("string"), S("value")),
	)
}

func TestAssociative(t *testing.T) {
	as := assert.New(t)
	m1 := getTestMap()

	as.Integer(3, api.Count(m1))

	nameKey := K("name")
	as.Equal(N("name"), nameKey.Name())

	nameValue, ok := m1.Get(nameKey)
	as.True(ok)
	as.String("Ale", nameValue)

	ageKey := K("age")
	ageValue, ok := m1.Get(ageKey)
	as.True(ok)
	as.Integer(99, ageValue)

	strValue, ok := m1.Get(S("string"))
	as.True(ok)
	as.String("value", strValue)

	r, ok := m1.Get(S("missing"))
	as.False(ok)
	as.Nil(r)
}

func TestAssociativeSequence(t *testing.T) {
	as := assert.New(t)
	m1 := getTestMap()

	first := m1.First()
	if e, ok := first.(api.Vector); ok {
		k, _ := e.ElementAt(0)
		v, _ := e.ElementAt(1)
		as.Equal(K("name"), k)
		as.String("Ale", v)
	} else {
		as.Fail("map.First() is not a vector")
	}

	rest := m1.Rest()
	as.String(`{:age 99, "string" "value"}`, rest)

}

func TestAssociativePrepend(t *testing.T) {
	as := assert.New(t)
	m1 := getTestMap()

	m2 := m1.Prepend(V(K("foo"), S("bar"))).(api.Associative)
	as.NotIdentical(m1, m2)

	r, ok := m2.Get(K("foo"))
	as.True(ok)
	as.String("bar", r)

	defer as.ExpectPanic(api.ExpectedPair)
	m2.Conjoin(F(99))
}
