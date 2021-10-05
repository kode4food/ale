package basic_test

import (
	"testing"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/types"
	"github.com/kode4food/ale/types/basic"
	"github.com/kode4food/ale/types/compound"
	"github.com/stretchr/testify/assert"
)

type testType struct{}

func TestBasic(t *testing.T) {
	as := assert.New(t)

	i99 := data.Integer(99).Type()
	i0 := data.Integer(0).Type()
	bTrue := data.True.Type()

	as.Equal("number", i99.Name())
	as.Equal("boolean", bTrue.Name())
	as.True(i99.Accepts(i0))
	as.False(i99.Accepts(bTrue))

	as.True(i99.Accepts(i99))
	as.False(i99.Accepts(basic.Any))
	as.False(i99.Accepts(compound.Union(basic.Symbol, basic.Pair)))

	as.False(i99.Accepts(testType{}))
}

func (testType) Name() string {
	return "test"
}

func (testType) Accepts(_ types.Type) bool {
	return false
}
