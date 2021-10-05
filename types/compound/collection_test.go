package compound_test

import (
	"testing"

	"github.com/kode4food/ale/types/basic"
	"github.com/kode4food/ale/types/compound"
	"github.com/stretchr/testify/assert"
)

func TestCollection(t *testing.T) {
	as := assert.New(t)

	c1 := compound.List(basic.Number)
	c2 := compound.Vector(basic.Number)

	as.Equal(basic.Number, c1.Element())
	as.True(c1.Accepts(c1))
	as.False(c1.Accepts(c2))

	as.True(basic.List.Accepts(c1))
	as.False(c1.Accepts(basic.List))
	as.False(basic.List.Accepts(c2))
	as.False(basic.Vector.Accepts(c1))
	as.True(basic.Vector.Accepts(c2))
	as.False(c2.Accepts(basic.Vector))
}
