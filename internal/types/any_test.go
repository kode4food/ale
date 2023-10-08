package types_test

import (
	"testing"

	"github.com/kode4food/ale/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestAny(t *testing.T) {
	as := assert.New(t)

	a := types.BasicAny
	as.Equal("any", a.Name())
	as.True(types.Accepts(a, types.BasicLambda))
	as.True(types.Accepts(a, types.BasicNumber))
	as.True(types.Accepts(a, types.BasicAny))
}
