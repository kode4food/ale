package data_test

import (
	"testing"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
)

func TestKeyword(t *testing.T) {
	as := assert.New(t)

	k1 := data.Keyword("hello")
	as.String("hello", k1.Name())
	as.String(":hello", k1)
}

func TestKeywordCaller(t *testing.T) {
	as := assert.New(t)

	m1 := O(C(data.Keyword("name"), S("Ale")))
	c1 := data.Keyword("name")
	as.String("Ale", c1.Call(m1))

	c2 := data.Keyword("missing")
	as.String("defaulted", c2.Call(m1, S("defaulted")))
}
