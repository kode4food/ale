package generate

import (
	"github.com/kode4food/ale/pkg/compiler/encoder"
	"github.com/kode4food/ale/pkg/runtime/isa"
)

// Branch constructs predicate, consequent, alternative branching
func Branch(e encoder.Encoder, predicate, consequent, alternative Builder) {
	thenLabel := e.NewLabel()
	endLabel := e.NewLabel()

	predicate(e)
	e.Emit(isa.CondJump, thenLabel)
	alternative(e)
	e.Emit(isa.Jump, endLabel)
	e.Emit(isa.Label, thenLabel)
	consequent(e)
	e.Emit(isa.Label, endLabel)
}
