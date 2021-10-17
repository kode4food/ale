package compound_test

import (
	"testing"

	"github.com/kode4food/ale/types"
	"github.com/kode4food/ale/types/basic"
	"github.com/kode4food/ale/types/compound"
	"github.com/stretchr/testify/assert"
)

func TestCons(t *testing.T) {
	as := assert.New(t)

	c1 := compound.Cons(basic.Null, basic.String)
	c2 := compound.Cons(basic.Number, basic.String)
	c3 := compound.Cons(basic.Null, basic.String)

	as.NotNil(types.Check(c1).Accepts(c1))
	as.Nil(types.Check(c1).Accepts(c2))
	as.NotNil(types.Check(c1).Accepts(c3))
}
