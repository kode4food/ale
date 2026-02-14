package types_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kode4food/ale/internal/types"
)

func TestAnyAccepts(t *testing.T) {
	as := assert.New(t)

	a := types.BasicAny
	as.Equal("any", a.Name())
	as.True(a.Accepts(types.BasicProcedure))
	as.True(a.Accepts(types.BasicNumber))
	as.True(a.Accepts(types.BasicAny))
}

func TestAnyEqual(t *testing.T) {
	as := assert.New(t)

	as.True(types.BasicAny.Equal(types.BasicAny))
	as.True(types.BasicAny.Equal(new(types.Any)))
	as.False(types.BasicAny.Equal(types.BasicBoolean))
}
