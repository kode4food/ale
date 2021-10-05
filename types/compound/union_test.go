package compound_test

import (
	"testing"

	"github.com/kode4food/ale/types/basic"
	"github.com/kode4food/ale/types/compound"
	"github.com/stretchr/testify/assert"
)

func TestUnion(t *testing.T) {
	as := assert.New(t)

	u1 := compound.Union(basic.Keyword, basic.Number)
	u2 := compound.Union(basic.List, basic.Vector)
	u3 := compound.Union(u1, u2)
	u4 := compound.Union(basic.Pair, basic.Symbol)
	u5 := compound.Union(basic.Number, basic.Bool, u1)

	as.Equal("keyword|number", u1.Name())
	as.Equal("keyword|list|number|vector", u3.Name())
	as.Equal("boolean|keyword|number", u5.Name())

	as.True(u1.Accepts(u1))
	as.True(u1.Accepts(basic.Keyword))
	as.True(u1.Accepts(basic.Number))
	as.False(u1.Accepts(basic.List))
	as.False(u1.Accepts(basic.Vector))

	as.True(u2.Accepts(basic.List))
	as.True(u2.Accepts(basic.Vector))
	as.False(u2.Accepts(basic.Keyword))
	as.False(u2.Accepts(basic.Number))

	as.True(u3.Accepts(basic.Keyword))
	as.True(u3.Accepts(basic.Number))
	as.True(u3.Accepts(basic.List))
	as.True(u3.Accepts(basic.Vector))
	as.False(u3.Accepts(basic.Symbol))

	as.True(u3.Accepts(u1))
	as.False(u4.Accepts(u3))

	as.False(basic.List.Accepts(u1))
}
