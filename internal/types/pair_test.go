package types_test

import (
	"testing"

	"github.com/kode4food/ale/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestCons(t *testing.T) {
	as := assert.New(t)

	c1 := types.MakeCons(types.BasicNull, types.BasicString)
	c2 := types.MakeCons(types.BasicNumber, types.BasicString)
	c3 := types.MakeCons(types.BasicNull, types.BasicString)

	as.True(types.Accepts(c1, c1))
	as.False(types.Accepts(c1, c2))
	as.True(types.Accepts(c1, c3))
}
