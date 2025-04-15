package types_test

import (
	"testing"

	"github.com/kode4food/ale/internal/types"
	"github.com/kode4food/ale/pkg/data"
	"github.com/stretchr/testify/assert"
)

type notABasic struct{}

func (n notABasic) Name() string                            { return "" }
func (n notABasic) Accepts(*types.Checker, types.Type) bool { return false }
func (n notABasic) Equal(types.Type) bool                   { return false }

func TestBasicAccepts(t *testing.T) {
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
	as.False(
		types.Accepts(i99, types.MakeUnion(types.BasicSymbol, types.BasicCons)),
	)
	as.False(types.Accepts(types.BasicBoolean, notABasic{}))
	as.True(types.Accepts(types.BasicAny, notABasic{}))
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
