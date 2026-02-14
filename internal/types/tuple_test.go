package types_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kode4food/ale/internal/types"
)

func TestTuple(t *testing.T) {
	as := assert.New(t)

	t1 := types.MakeTuple(types.BasicKeyword, types.BasicNumber)
	t2 := types.MakeTuple(types.BasicNumber, types.BasicKeyword)
	t3 := types.MakeTuple(types.BasicKeyword, types.BasicNumber)
	t4 := types.MakeTuple()

	as.Equal("tuple(keyword,number)", t1.Name())

	as.True(t1.Accepts(t1))
	as.True(t1.Accepts(t3))
	as.False(t2.Accepts(t1))
	as.False(t2.Accepts(t4))
	as.False(t4.Accepts(t1))

	as.False(t1.Accepts(types.BasicNull))
}
