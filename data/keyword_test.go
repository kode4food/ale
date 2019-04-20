package data_test

import (
	"testing"

	"gitlab.com/kode4food/ale/internal/assert"
	. "gitlab.com/kode4food/ale/internal/assert/helpers"
)

func TestKeyword(t *testing.T) {
	as := assert.New(t)

	k1 := K("hello")
	as.String("hello", k1.Name())
	as.String(":hello", k1)
}

func TestKeywordCaller(t *testing.T) {
	as := assert.New(t)

	m1 := getTestMap()
	c1 := K("name").Caller()
	as.String("Ale", c1(m1))

	c2 := K("missing").Caller()
	as.String("defaulted", c2(m1, S("defaulted")))
}
