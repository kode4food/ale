package assert_test

import (
	"testing"

	"github.com/kode4food/ale/internal/assert"
	"github.com/kode4food/ale/internal/runtime/isa"
	"github.com/kode4food/ale/pkg/data"
)

func TestEncodesAs(t *testing.T) {
	as := assert.New(t)
	as.EncodesAs(isa.Instructions{
		isa.PosInt.New(2),
	}, `2`)
}

func TestGetRootSymbol(t *testing.T) {
	as := assert.New(t)

	e := assert.GetTestEncoder()
	v1 := assert.GetRootSymbol(e, "true")
	v2 := assert.GetRootSymbol(e, "false")
	as.Equal(data.True, v1)
	as.Equal(data.False, v2)
}
