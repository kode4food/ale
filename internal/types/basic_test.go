package types_test

import (
	"testing"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestBasicAccepts(t *testing.T) {
	as := assert.New(t)

	i99 := data.Integer(99).Type()
	i0 := data.Integer(0).Type()
	bTrue := data.True.Type()
	bFalse := data.False.Type()

	as.Equal("number(99)", i99.Name())
	as.Equal("boolean(#t)", bTrue.Name())

	as.True(types.BasicNumber.Accepts(i99))
	as.True(types.BasicNumber.Accepts(i0))
	as.False(i99.Accepts(types.BasicNumber))
	as.False(i0.Accepts(types.BasicNumber))
	as.False(i99.Accepts(i0))
	as.False(i99.Accepts(bTrue))

	as.True(i99.Accepts(i99))
	as.False(i99.Accepts(types.BasicAny))
	as.True(types.BasicAny.Accepts(i99))
	as.False(
		i99.Accepts(types.MakeUnion(types.BasicSymbol, types.BasicCons)),
	)

	as.True(types.BasicBoolean.Accepts(bTrue))
	as.True(types.BasicBoolean.Accepts(bFalse))
	as.False(bTrue.Accepts(types.BasicBoolean))
	as.False(bFalse.Accepts(types.BasicBoolean))
}

func TestBasicEqual(t *testing.T) {
	as := assert.New(t)

	i99 := data.Integer(99).Type()
	i0 := data.Integer(0).Type()
	bTrue := data.True.Type()

	as.True(i99.Equal(data.Integer(99).Type()))
	as.False(i99.Equal(i0))
	as.False(i99.Equal(bTrue))
}
