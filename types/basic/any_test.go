package basic_test

import (
	"testing"

	"github.com/kode4food/ale/types"
	"github.com/kode4food/ale/types/basic"
	"github.com/stretchr/testify/assert"
)

func TestAny(t *testing.T) {
	as := assert.New(t)

	as.Equal("any", basic.Any.Name())
	as.NotNil(types.Check(basic.Any).Accepts(basic.Lambda))
	as.NotNil(types.Check(basic.Any).Accepts(basic.Number))
	as.NotNil(types.Check(basic.Any).Accepts(basic.Any))
}
