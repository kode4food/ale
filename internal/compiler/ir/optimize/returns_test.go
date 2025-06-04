package optimize_test

import (
	"testing"

	"github.com/kode4food/ale/internal/assert"
	"github.com/kode4food/ale/internal/compiler/encoder"
	"github.com/kode4food/ale/internal/compiler/generate"
	"github.com/kode4food/ale/internal/compiler/ir/optimize"
	"github.com/kode4food/ale/internal/runtime/isa"
)

func TestSplitReturns(t *testing.T) {
	as := assert.New(t)

	e1 := assert.GetTestEncoder()
	as.NoError(generate.Branch(e1,
		func(encoder.Encoder) error { e1.Emit(isa.True); return nil },
		func(encoder.Encoder) error { e1.Emit(isa.PosInt, 1); return nil },
		func(encoder.Encoder) error { e1.Emit(isa.Zero); return nil },
	))
	e1.Emit(isa.Return)

	as.Instructions(isa.Instructions{
		isa.True.New(),
		isa.CondJump.New(0),
		isa.Zero.New(),
		isa.Return.New(),
		isa.Jump.New(1),
		isa.Label.New(0),
		isa.PosInt.New(1),
		isa.Return.New(),
		isa.Label.New(1),
	}, optimize.Encoded(e1.Encode()).Code)
}
