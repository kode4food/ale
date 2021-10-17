package compound_test

import (
	"testing"

	"github.com/kode4food/ale/types"
	"github.com/kode4food/ale/types/basic"
	"github.com/kode4food/ale/types/compound"
	"github.com/stretchr/testify/assert"
)

func TestTuple(t *testing.T) {
	as := assert.New(t)

	t1 := compound.Tuple(basic.Keyword, basic.Number)
	t2 := compound.Tuple(basic.Number, basic.Keyword)
	t3 := compound.Tuple(basic.Keyword, basic.Number)
	t4 := compound.Tuple()

	as.Equal("tuple(keyword,number)", t1.Name())

	as.NotNil(types.Check(t1).Accepts(t1))
	as.NotNil(types.Check(t1).Accepts(t3))
	as.Nil(types.Check(t2).Accepts(t1))
	as.Nil(types.Check(t2).Accepts(t4))
	as.Nil(types.Check(t4).Accepts(t1))

	as.Nil(types.Check(t1).Accepts(basic.Null))
}
