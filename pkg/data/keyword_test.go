package data_test

import (
	"testing"

	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/internal/types"
	"github.com/kode4food/ale/pkg/data"
)

func TestKeyword(t *testing.T) {
	as := assert.New(t)

	k1 := data.Keyword("hello")
	as.String("hello", k1.Name())
	as.String(":hello", k1)

	as.True(types.BasicKeyword.Equal(k1.Type()))
}

func TestKeywordCaller(t *testing.T) {
	as := assert.New(t)

	m1 := O(C(data.Keyword("name"), S("Ale")))
	c1 := data.Keyword("name")
	as.String("Ale", c1.Call(m1))

	c2 := data.Keyword("missing")
	as.String("defaulted", c2.Call(m1, S("defaulted")))
}
