package types_test

import (
	"testing"

	"github.com/kode4food/ale/types"
	"github.com/stretchr/testify/assert"
)

func TestCons(t *testing.T) {
	as := assert.New(t)

	c1 := types.Cons(types.Null, types.String)
	c2 := types.Cons(types.Number, types.String)
	c3 := types.Cons(types.Null, types.String)

	as.True(types.Accepts(c1, c1))
	as.False(types.Accepts(c1, c2))
	as.True(types.Accepts(c1, c3))
}
