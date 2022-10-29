package types_test

import (
	"testing"

	"github.com/kode4food/ale/types"
	"github.com/stretchr/testify/assert"
)

func TestList(t *testing.T) {
	as := assert.New(t)

	v1 := types.ListOf(types.Number)
	v2 := types.ListOf(types.String)
	v3 := types.ListOf(types.Number)

	as.Equal("list(number)", v1.Name())

	as.True(types.Accepts(v1, v1))
	as.False(types.Accepts(v1, v2))
	as.True(types.Accepts(v3, v1))

	as.False(types.Accepts(types.AnyVector, v1))
	as.True(types.Accepts(types.AnyList, v1))
	as.False(types.Accepts(v1, types.AnyList))
}

func TestVector(t *testing.T) {
	as := assert.New(t)

	v1 := types.VectorOf(types.Number)
	v2 := types.VectorOf(types.String)
	v3 := types.VectorOf(types.Number)

	as.Equal("vector(number)", v1.Name())

	as.True(types.Accepts(v1, v1))
	as.False(types.Accepts(v1, v2))
	as.True(types.Accepts(v3, v1))

	as.False(types.Accepts(types.AnyList, v1))
	as.True(types.Accepts(types.AnyVector, v1))
	as.False(types.Accepts(v1, types.AnyVector))
}
