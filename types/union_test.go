package types_test

import (
	"testing"

	"github.com/kode4food/ale/types"
	"github.com/stretchr/testify/assert"
)

func TestUnion(t *testing.T) {
	as := assert.New(t)

	u1 := types.Union(types.Keyword, types.Number)
	u2 := types.Union(types.AnyList, types.AnyVector)
	u3 := types.Union(u1, u2)
	u4 := types.Union(types.AnyCons, types.Symbol)
	u5 := types.Union(types.Number, types.Bool, u1)
	u6 := types.Union(u1, u2, types.Any, u5)
	u7 := types.Union(types.Number, types.Keyword)

	as.Equal("union(keyword,number)", u1.Name())
	as.Equal("union(keyword,list,number,vector)", u3.Name())
	as.Equal("union(boolean,keyword,number)", u5.Name())
	as.Equal("any", u6.Name())

	as.True(types.Accepts(u1, u1))
	as.True(types.Accepts(u1, types.Keyword))
	as.True(types.Accepts(u1, types.Number))
	as.False(types.Accepts(u1, types.AnyList))
	as.False(types.Accepts(u1, types.AnyVector))

	as.True(types.Accepts(u2, types.AnyList))
	as.True(types.Accepts(u2, types.AnyVector))
	as.False(types.Accepts(u2, types.Keyword))
	as.False(types.Accepts(u2, types.Number))

	as.True(types.Accepts(u3, types.Keyword))
	as.True(types.Accepts(u3, types.Number))
	as.True(types.Accepts(u3, types.AnyList))
	as.True(types.Accepts(u3, types.AnyVector))
	as.False(types.Accepts(u3, types.Symbol))

	as.True(types.Accepts(u3, u1))
	as.False(types.Accepts(u4, u3))

	as.False(types.Accepts(types.AnyList, u1))

	as.True(types.Accepts(u7, u1))
	as.True(types.Accepts(u1, u7))

	_, ok := u6.(types.AnyType)
	as.True(ok)
}
