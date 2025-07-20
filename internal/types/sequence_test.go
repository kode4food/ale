package types_test

import (
	"testing"

	"github.com/kode4food/ale/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestList(t *testing.T) {
	as := assert.New(t)

	v1 := types.MakeListOf(types.BasicNumber)
	v2 := types.MakeListOf(types.BasicString)
	v3 := types.MakeListOf(types.BasicNumber)

	as.Equal("list(number)", v1.Name())

	as.True(v1.Accepts(v1))
	as.False(v1.Accepts(v2))
	as.True(v3.Accepts(v1))

	as.False(types.BasicVector.Accepts(v1))
	as.True(types.BasicList.Accepts(v1))
	as.False(v1.Accepts(types.BasicList))
}

func TestVector(t *testing.T) {
	as := assert.New(t)

	v1 := types.MakeVectorOf(types.BasicNumber)
	v2 := types.MakeVectorOf(types.BasicString)
	v3 := types.MakeVectorOf(types.BasicNumber)

	as.Equal("vector(number)", v1.Name())

	as.True(v1.Accepts(v1))
	as.False(v1.Accepts(v2))
	as.True(v3.Accepts(v1))

	as.False(types.BasicList.Accepts(v1))
	as.True(types.BasicVector.Accepts(v1))
	as.False(v1.Accepts(types.BasicVector))
}
