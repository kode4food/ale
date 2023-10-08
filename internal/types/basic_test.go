package types_test

import (
	"testing"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestBasic(t *testing.T) {
	as := assert.New(t)

	i99 := data.Integer(99).Type()
	i0 := data.Integer(0).Type()
	bTrue := data.True.Type()

	as.Equal("number", i99.Name())
	as.Equal("boolean", bTrue.Name())
	as.True(types.Accepts(i99, i0))
	as.False(types.Accepts(i99, bTrue))

	as.True(types.Accepts(i99, i99))
	as.False(types.Accepts(i99, types.BasicAny))
	as.False(types.Accepts(i99, types.MakeUnion(types.BasicSymbol, types.BasicCons)))
}
