package optimize_test

import (
	"testing"

	"github.com/kode4food/ale/compiler/ir/optimize"
	"github.com/kode4food/ale/internal/assert"
	"github.com/kode4food/ale/runtime/isa"
)

func TestUntruthy(t *testing.T) {
	as := assert.New(t)

	e1 := assert.GetTestEncoder()
	e1.Emit(isa.False)
	e1.Emit(isa.MakeTruthy)
	e1.Emit(isa.Return)

	as.Instructions(isa.Instructions{
		isa.New(isa.RetFalse),
	}, optimize.Instructions(e1.Code()))
}
