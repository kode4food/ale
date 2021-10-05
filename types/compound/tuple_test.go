package compound_test

import (
	"testing"

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

	as.Equal("tuple", t1.Name())

	as.True(t1.Accepts(t1))
	as.True(t1.Accepts(t3))
	as.False(t2.Accepts(t1))
	as.False(t2.Accepts(t4))
	as.False(t4.Accepts(t1))

	as.False(t1.Accepts(basic.Null))
}

func TestVectorTuple(t *testing.T) {
	as := assert.New(t)

	t1 := compound.VectorTuple(basic.Number, basic.Keyword)
	t2 := compound.Tuple(basic.Number, basic.Keyword)
	t3 := compound.ListTuple(basic.Number, basic.Keyword)

	as.True(t2.Accepts(t1))
	as.False(t1.Accepts(t2))
	as.False(t1.Accepts(t3))
}

func TestListTuple(t *testing.T) {
	as := assert.New(t)

	t1 := compound.ListTuple(basic.Number, basic.Keyword)
	t2 := compound.Tuple(basic.Number, basic.Keyword)
	t3 := compound.VectorTuple(basic.Number, basic.Keyword)

	as.True(t2.Accepts(t1))
	as.False(t1.Accepts(t2))
	as.False(t1.Accepts(t3))
}
