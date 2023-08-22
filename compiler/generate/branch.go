package generate

import (
	"github.com/kode4food/ale/compiler/encoder"
	"github.com/kode4food/ale/runtime/isa"
)

// Branch constructs predicate, consequent, alternative branching
func Branch(e encoder.Encoder, predicate, consequent, alternative func()) {
	thenLabel := e.NewLabel()
	endLabel := e.NewLabel()

	predicate()
	e.Emit(isa.CondJump, thenLabel)
	alternative()
	e.Emit(isa.Jump, endLabel)
	e.Emit(isa.Label, thenLabel)
	consequent()
	e.Emit(isa.Label, endLabel)
}
