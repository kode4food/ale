package types_test

import (
	"testing"

	"github.com/kode4food/ale/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestConsAccepts(t *testing.T) {
	as := assert.New(t)

	c1 := types.MakeCons(types.BasicNull, types.BasicString)
	c2 := types.MakeCons(types.BasicNumber, types.BasicString)
	c3 := types.MakeCons(types.BasicNull, types.BasicString)

	as.True(c1.Accepts(c1))
	as.False(c1.Accepts(c2))
	as.True(c1.Accepts(c3))
}

func TestConsEqual(t *testing.T) {
	as := assert.New(t)

	c1 := types.MakeCons(types.BasicNull, types.BasicString)
	c2 := types.MakeCons(types.BasicNumber, types.BasicString)
	c3 := types.MakeCons(types.BasicNull, types.BasicString)

	as.True(c1.Equal(c1))
	as.False(c1.Equal(c2))
	as.True(c1.Equal(c3))

	as.False(types.BasicObject.Equal(c1))
	as.False(c1.Equal(types.BasicObject))
	as.False(c1.Equal(types.BasicBoolean))
}
