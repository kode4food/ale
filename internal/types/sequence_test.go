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

	as.True(types.Accepts(v1, v1))
	as.False(types.Accepts(v1, v2))
	as.True(types.Accepts(v3, v1))

	as.False(types.Accepts(types.BasicVector, v1))
	as.True(types.Accepts(types.BasicList, v1))
	as.False(types.Accepts(v1, types.BasicList))
}

func TestVector(t *testing.T) {
	as := assert.New(t)

	v1 := types.MakeVectorOf(types.BasicNumber)
	v2 := types.MakeVectorOf(types.BasicString)
	v3 := types.MakeVectorOf(types.BasicNumber)

	as.Equal("vector(number)", v1.Name())

	as.True(types.Accepts(v1, v1))
	as.False(types.Accepts(v1, v2))
	as.True(types.Accepts(v3, v1))

	as.False(types.Accepts(types.BasicList, v1))
	as.True(types.Accepts(types.BasicVector, v1))
	as.False(types.Accepts(v1, types.BasicVector))
}
