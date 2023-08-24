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
		func(encoder.Encoder) { e.Emit(isa.One) },
		func(encoder.Encoder) { e.Emit(isa.Zero) },
	)
	e.Emit(isa.Return)

	as.Instructions(
		isa.Instructions{
			isa.New(isa.True),
			isa.New(isa.CondJump, 0),
			isa.New(isa.Zero),
			isa.New(isa.Jump, 1),
			isa.New(isa.Label, 0),
			isa.New(isa.One),
			isa.New(isa.Label, 1),
			isa.New(isa.Return),
		},
		e.Code(),
	)
}
