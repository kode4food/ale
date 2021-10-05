package types_test

import (
	"testing"

	"github.com/kode4food/ale/types"
	"github.com/kode4food/ale/types/basic"
	"github.com/stretchr/testify/assert"
)

func TestAny(t *testing.T) {
	as := assert.New(t)

	as.Equal("any", types.Any.Name())
	as.True(types.Any.Accepts(basic.Lambda))
	as.True(types.Any.Accepts(basic.Number))
	as.True(types.Any.Accepts(types.Any))
}
