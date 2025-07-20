package types_test

import (
	"testing"

	"github.com/kode4food/ale"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/types"
	"github.com/stretchr/testify/assert"
)

type notABasic struct{}

func (n notABasic) Name() string          { return "" }
func (n notABasic) Accepts(ale.Type) bool { return false }
func (n notABasic) Equal(ale.Type) bool   { return false }

func TestBasicAccepts(t *testing.T) {
	as := assert.New(t)

	i99 := data.Integer(99).Type()
	i0 := data.Integer(0).Type()
	bTrue := data.True.Type()

	as.Equal("number", i99.Name())
	as.Equal("boolean", bTrue.Name())
	as.True(i99.Accepts(i0))
	as.False(i99.Accepts(bTrue))

	as.True(i99.Accepts(i99))
	as.False(i99.Accepts(types.BasicAny))
	as.False(
		i99.Accepts(types.MakeUnion(types.BasicSymbol, types.BasicCons)),
	)
	as.False(types.BasicBoolean.Accepts(notABasic{}))
	as.True(types.BasicAny.Accepts(notABasic{}))
}

func TestBasicEqual(t *testing.T) {
	as := assert.New(t)

	i99 := data.Integer(99).Type()
	i0 := data.Integer(0).Type()
	bTrue := data.True.Type()

	as.True(i99.Equal(i0))
	as.False(i99.Equal(bTrue))
	as.False(types.BasicBoolean.Equal(notABasic{}))
}
