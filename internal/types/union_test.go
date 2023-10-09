package types_test

import (
	"testing"

	"github.com/kode4food/ale/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestUnionAccepts(t *testing.T) {
	as := assert.New(t)

	u1 := types.MakeUnion(types.BasicKeyword, types.BasicNumber)
	u2 := types.MakeUnion(types.BasicList, types.BasicVector)
	u3 := types.MakeUnion(u1, u2)
	u4 := types.MakeUnion(types.BasicCons, types.BasicSymbol)
	u5 := types.MakeUnion(types.BasicNumber, types.BasicBoolean, u1)
	u6 := types.MakeUnion(u1, u2, types.BasicAny, u5)
	u7 := types.MakeUnion(types.BasicNumber, types.BasicKeyword)

	as.Equal("union(keyword,number)", u1.Name())
	as.Equal("union(keyword,list,number,vector)", u3.Name())
	as.Equal("union(boolean,keyword,number)", u5.Name())
	as.Equal("any", u6.Name())

	as.True(types.Accepts(u1, u1))
	as.True(types.Accepts(u1, types.BasicKeyword))
	as.True(types.Accepts(u1, types.BasicNumber))
	as.False(types.Accepts(u1, types.BasicList))
	as.False(types.Accepts(u1, types.BasicVector))

	as.True(types.Accepts(u2, types.BasicList))
	as.True(types.Accepts(u2, types.BasicVector))
	as.False(types.Accepts(u2, types.BasicKeyword))
	as.False(types.Accepts(u2, types.BasicNumber))

	as.True(types.Accepts(u3, types.BasicKeyword))
	as.True(types.Accepts(u3, types.BasicNumber))
	as.True(types.Accepts(u3, types.BasicList))
	as.True(types.Accepts(u3, types.BasicVector))
	as.False(types.Accepts(u3, types.BasicSymbol))

	as.True(types.Accepts(u3, u1))
	as.False(types.Accepts(u4, u3))

	as.False(types.Accepts(types.BasicList, u1))

	as.True(types.Accepts(u7, u1))
	as.True(types.Accepts(u1, u7))

	_, ok := u6.(*types.Any)
	as.True(ok)
}

func TestUnionEqual(t *testing.T) {
	as := assert.New(t)

	u1 := types.MakeUnion(types.BasicKeyword, types.BasicNumber)
	u2 := types.MakeUnion(types.BasicList, types.BasicVector)
	u3 := types.MakeUnion(u1, u2)
	u4 := types.MakeUnion(
		types.BasicList, types.BasicVector,
		types.BasicKeyword, types.BasicNumber,
	)
	u5 := types.MakeUnion(types.BasicNumber, types.BasicKeyword)

	as.True(u1.Equal(u1))
	as.False(u1.Equal(types.BasicKeyword))
	as.True(u1.Equal(u5))
	as.False(types.BasicNumber.Equal(u1))
	as.True(u3.Equal(u4))
}
