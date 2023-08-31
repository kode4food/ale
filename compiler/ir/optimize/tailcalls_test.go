package optimize_test

import (
	"testing"

	"github.com/kode4food/ale/compiler/generate"
	"github.com/kode4food/ale/compiler/ir/optimize"
	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/runtime/isa"
)

func TestTailCalls(t *testing.T) {
	as := assert.New(t)

	e1 := assert.GetTestEncoder()
	generate.Call(e1, L(assert.GetRootSymbol(e1, "+"), I(1), I(2)))
	e1.Emit(isa.Return)

	as.Instructions(isa.Instructions{
		isa.New(isa.PosInt, 2),
		isa.New(isa.PosInt, 1),
		isa.New(isa.Const, 0),
		isa.New(isa.TailCall, 2),
	}, optimize.Instructions(e1.Code()))
}
