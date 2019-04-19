package generate

import (
	"gitlab.com/kode4food/ale/compiler/encoder"
	"gitlab.com/kode4food/ale/runtime/isa"
)

// Builder is a function that a composer (ex: Branch) will invoke
type Builder func()

// Branch constructs conditional branching
func Branch(e encoder.Type, cond, thenBranch, elseBranch Builder) {
	thenLabel := e.NewLabel()
	endLabel := e.NewLabel()

	cond()
	e.Emit(isa.CondJump, thenLabel)
	elseBranch()
	e.Emit(isa.Jump, endLabel)
	thenLabel.DropAnchor()
	thenBranch()
	endLabel.DropAnchor()
}
