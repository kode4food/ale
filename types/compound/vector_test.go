package compound_test

import (
	"testing"

	"github.com/kode4food/ale/types"
	"github.com/kode4food/ale/types/basic"
	"github.com/kode4food/ale/types/compound"
	"github.com/stretchr/testify/assert"
)

func TestVector(t *testing.T) {
	as := assert.New(t)

	v1 := compound.Vector(basic.Number)
	v2 := compound.Vector(basic.String)
	v3 := compound.Vector(basic.Number)

	as.Equal("vector(number)", v1.Name())

	as.NotNil(types.Check(v1).Accepts(v1))
	as.Nil(types.Check(v1).Accepts(v2))
	as.NotNil(types.Check(v3).Accepts(v1))

	as.Nil(types.Check(basic.List).Accepts(v1))
	as.NotNil(types.Check(basic.Vector).Accepts(v1))
	as.Nil(types.Check(v1).Accepts(basic.Vector))
}
