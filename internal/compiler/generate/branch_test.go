package generate_test

import (
	"testing"

	"github.com/kode4food/ale/internal/assert"
	"github.com/kode4food/ale/internal/compiler/encoder"
	"github.com/kode4food/ale/internal/compiler/generate"
	"github.com/kode4food/ale/internal/runtime/isa"
)

func TestBranch(t *testing.T) {
	as := assert.New(t)

	e := assert.GetTestEncoder()
	as.Nil(generate.Branch(e,
		func(encoder.Encoder) error { e.Emit(isa.True); return nil },
		func(encoder.Encoder) error { e.Emit(isa.PosInt, 1); return nil },
		func(encoder.Encoder) error { e.Emit(isa.Zero); return nil },
	))
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
		e.Encode().Code,
	)
}
