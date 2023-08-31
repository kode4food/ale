package optimize_test

import (
	"testing"

	"github.com/kode4food/ale/compiler/encoder"
	"github.com/kode4food/ale/compiler/generate"
	"github.com/kode4food/ale/compiler/ir/optimize"
	"github.com/kode4food/ale/internal/assert"
	"github.com/kode4food/ale/runtime/isa"
)

func TestSplitReturns(t *testing.T) {
	as := assert.New(t)

	e1 := assert.GetTestEncoder()
	generate.Branch(e1,
		func(encoder.Encoder) { e1.Emit(isa.True) },
		func(encoder.Encoder) { e1.Emit(isa.PosInt, 1) },
		func(encoder.Encoder) { e1.Emit(isa.Zero) },
	)
	e1.Emit(isa.Return)

	as.Instructions(isa.Instructions{
		isa.New(isa.True),
		isa.New(isa.CondJump, 0),
		isa.New(isa.Zero),
		isa.New(isa.Return),
		isa.New(isa.Jump, 1),
		isa.New(isa.Label, 0),
		isa.New(isa.PosInt, 1),
		isa.New(isa.Return),
		isa.New(isa.Label, 1),
	}, optimize.Instructions(e1.Code()))
}
