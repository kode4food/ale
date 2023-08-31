package generate_test

import (
	"testing"

	"github.com/kode4food/ale/compiler/encoder"
	"github.com/kode4food/ale/compiler/generate"
	"github.com/kode4food/ale/internal/assert"
	"github.com/kode4food/ale/runtime/isa"
)

func TestBranch(t *testing.T) {
	as := assert.New(t)

	e := assert.GetTestEncoder()
	generate.Branch(e,
		func(encoder.Encoder) { e.Emit(isa.True) },
		func(encoder.Encoder) { e.Emit(isa.PosInt, 1) },
		func(encoder.Encoder) { e.Emit(isa.Zero) },
	)
	e.Emit(isa.Return)

	as.Instructions(
		isa.Instructions{
			isa.True.New(),
			isa.CondJump.New(0),
			isa.Zero.New(),
			isa.Jump.New(1),
			isa.Label.New(0),
			isa.PosInt.New(1),
			isa.Label.New(1),
			isa.Return.New(),
		},
		e.Code(),
	)
}
