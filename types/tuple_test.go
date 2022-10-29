package types_test

import (
	"testing"

	"github.com/kode4food/ale/types"
	"github.com/stretchr/testify/assert"
)

func TestTuple(t *testing.T) {
	as := assert.New(t)

	t1 := types.Tuple(types.Keyword, types.Number)
	t2 := types.Tuple(types.Number, types.Keyword)
	t3 := types.Tuple(types.Keyword, types.Number)
	t4 := types.Tuple()

	as.Equal("tuple(keyword,number)", t1.Name())

	as.True(types.Accepts(t1, t1))
	as.True(types.Accepts(t1, t3))
	as.False(types.Accepts(t2, t1))
	as.False(types.Accepts(t2, t4))
	as.False(types.Accepts(t4, t1))

	as.False(types.Accepts(t1, types.Null))
}
