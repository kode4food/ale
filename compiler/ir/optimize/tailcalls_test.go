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
		isa.PosInt.New(2),
		isa.PosInt.New(1),
		isa.Const.New(0),
		isa.TailCall.New(2),
	}, optimize.FromEncoded(e1.Encode()))
}
