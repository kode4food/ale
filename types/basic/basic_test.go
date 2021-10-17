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
	as.NotNil(types.Check(i99).Accepts(i0))
	as.Nil(types.Check(i99).Accepts(bTrue))

	as.NotNil(types.Check(i99).Accepts(i99))
	as.Nil(types.Check(i99).Accepts(basic.Any))
	as.Nil(types.Check(i99).Accepts(compound.Union(basic.Symbol, basic.Cons)))

	as.Nil(types.Check(i99).Accepts(testType{}))
}

func (testType) Name() string {
	return "test"
}

func (testType) Accepts(_ types.Checker, _ types.Type) bool {
	return false
}
