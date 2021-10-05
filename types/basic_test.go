package types_test

import (
	"testing"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/types"
	"github.com/stretchr/testify/assert"
)

func TestBasic(t *testing.T) {
	as := assert.New(t)

	i99 := data.Integer(99).Type()
	i0 := data.Integer(0).Type()
	bTrue := data.True.Type()

	as.Equal("number", i99.Name())
	as.Equal("boolean", bTrue.Name())
	as.True(i99.Accepts(i0))
	as.False(i99.Accepts(bTrue))

	as.False(i99.Accepts(types.Any))
}
