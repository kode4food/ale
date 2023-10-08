package types_test

import (
	"testing"

	"github.com/kode4food/ale/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestTuple(t *testing.T) {
	as := assert.New(t)

	t1 := types.MakeTuple(types.BasicKeyword, types.BasicNumber)
	t2 := types.MakeTuple(types.BasicNumber, types.BasicKeyword)
	t3 := types.MakeTuple(types.BasicKeyword, types.BasicNumber)
	t4 := types.MakeTuple()

	as.Equal("tuple(keyword,number)", t1.Name())

	as.True(types.Accepts(t1, t1))
	as.True(types.Accepts(t1, t3))
	as.False(types.Accepts(t2, t1))
	as.False(types.Accepts(t2, t4))
	as.False(types.Accepts(t4, t1))

	as.False(types.Accepts(t1, types.BasicNull))
}
