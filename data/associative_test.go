package data_test

import (
	"testing"

	"gitlab.com/kode4food/ale/data"
	"gitlab.com/kode4food/ale/internal/assert"
	. "gitlab.com/kode4food/ale/internal/assert/helpers"
)

func getTestMap() data.Associative {
	return A(
		K("name"), S("Ale"),
		K("age"), I(99),
		S("string"), S("value"),
	)
}

func TestAssociative(t *testing.T) {
	as := assert.New(t)
	m1 := getTestMap()

	as.Number(3, m1.Count())

	nameKey := K("name")
	as.Equal(N("name"), nameKey.Name())

	nameValue, ok := m1.Get(nameKey)
	as.True(ok)
	as.String("Ale", nameValue)

	ageKey := K("age")
	ageValue, ok := m1.Get(ageKey)
	as.True(ok)
	as.Number(99, ageValue)

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
	if e, ok := first.(data.Vector); ok {
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

func TestAssociativeSplit(t *testing.T) {
	as := assert.New(t)
	m1 := getTestMap()

	f, r, ok := m1.Split()
	as.True(ok)
	as.False(r.IsEmpty())
	v, ok := f.(data.Vector)
	as.True(ok)
	as.Number(2, v.Count())

	f, ok = v.ElementAt(0)
	as.True(ok)
	as.String(":name", f)

	f, ok = v.ElementAt(1)
	as.True(ok)
	as.String("Ale", f)
}

func TestAssociativePrepend(t *testing.T) {
	as := assert.New(t)
	m1 := getTestMap()

	m2 := m1.Prepend(V(K("foo"), S("bar"))).(data.Associative)
	as.NotIdentical(m1, m2)

	r, ok := m2.Get(K("foo"))
	as.True(ok)
	as.String("bar", r)
}

func TestEmptyAssociative(t *testing.T) {
	as := assert.New(t)

	m1 := data.EmptyAssociative
	as.Nil(m1.First())
	as.True(m1.IsEmpty())
	as.True(m1.Rest().IsEmpty())

	f, r, ok := m1.Split()
	as.Nil(f)
	as.True(r.IsEmpty())
	as.False(ok)
}

func TestAssociativeCaller(t *testing.T) {
	as := assert.New(t)

	m1 := getTestMap()
	c := m1.Caller()
	as.String("Ale", c(K("name")))
	as.Nil(c(K("unknown")))
	as.String("defaulted", c(K("unknown"), S("defaulted")))
}
