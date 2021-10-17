package compound_test

import (
	"testing"

	"github.com/kode4food/ale/types"
	"github.com/kode4food/ale/types/basic"
	"github.com/kode4food/ale/types/compound"
	"github.com/stretchr/testify/assert"
)

func TestList(t *testing.T) {
	as := assert.New(t)

	v1 := compound.List(basic.Number)
	v2 := compound.List(basic.String)
	v3 := compound.List(basic.Number)

	as.Equal("list(number)", v1.Name())

	as.NotNil(types.Check(v1).Accepts(v1))
	as.Nil(types.Check(v1).Accepts(v2))
	as.NotNil(types.Check(v3).Accepts(v1))

	as.Nil(types.Check(basic.Vector).Accepts(v1))
	as.NotNil(types.Check(basic.List).Accepts(v1))
	as.Nil(types.Check(v1).Accepts(basic.List))
}
