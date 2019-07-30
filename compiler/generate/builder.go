package generate

import (
	"github.com/kode4food/ale/compiler/encoder"
	"github.com/kode4food/ale/runtime/isa"
)

// Builder is a function that a composer (ex: Branch) will invoke
type Builder func()

// Branch constructs conditional branching
func Branch(e encoder.Type, predicate, consequent, alternative Builder) {
	thenLabel := e.NewLabel()
	endLabel := e.NewLabel()

	predicate()
	e.Emit(isa.CondJump, thenLabel)
	alternative()
	e.Emit(isa.Jump, endLabel)
	thenLabel.DropAnchor()
	consequent()
	endLabel.DropAnchor()
}
