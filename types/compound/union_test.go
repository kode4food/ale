package compound_test

import (
	"testing"

	"github.com/kode4food/ale/types"
	"github.com/kode4food/ale/types/basic"
	"github.com/kode4food/ale/types/compound"
	"github.com/stretchr/testify/assert"
)

func TestUnion(t *testing.T) {
	as := assert.New(t)

	u1 := compound.Union(basic.Keyword, basic.Number)
	u2 := compound.Union(basic.List, basic.Vector)
	u3 := compound.Union(u1, u2)
	u4 := compound.Union(basic.Cons, basic.Symbol)
	u5 := compound.Union(basic.Number, basic.Bool, u1)
	u6 := compound.Union(u1, u2, basic.Any, u5)

	as.Equal("union(keyword,number)", u1.Name())
	as.Equal("union(keyword,list,number,vector)", u3.Name())
	as.Equal("union(boolean,keyword,number)", u5.Name())
	as.Equal("any", u6.Name())

	as.NotNil(types.Check(u1).Accepts(u1))
	as.NotNil(types.Check(u1).Accepts(basic.Keyword))
	as.NotNil(types.Check(u1).Accepts(basic.Number))
	as.Nil(types.Check(u1).Accepts(basic.List))
	as.Nil(types.Check(u1).Accepts(basic.Vector))

	as.NotNil(types.Check(u2).Accepts(basic.List))
	as.NotNil(types.Check(u2).Accepts(basic.Vector))
	as.Nil(types.Check(u2).Accepts(basic.Keyword))
	as.Nil(types.Check(u2).Accepts(basic.Number))

	as.NotNil(types.Check(u3).Accepts(basic.Keyword))
	as.NotNil(types.Check(u3).Accepts(basic.Number))
	as.NotNil(types.Check(u3).Accepts(basic.List))
	as.NotNil(types.Check(u3).Accepts(basic.Vector))
	as.Nil(types.Check(u3).Accepts(basic.Symbol))

	as.NotNil(types.Check(u3).Accepts(u1))
	as.Nil(types.Check(u4).Accepts(u3))

	as.Nil(types.Check(basic.List).Accepts(u1))

	_, ok := u6.(basic.AnyType)
	as.True(ok)
}
