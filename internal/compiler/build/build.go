package build

import (
	"gitlab.com/kode4food/ale/internal/compiler/encoder"
	"gitlab.com/kode4food/ale/internal/runtime/isa"
)

// Builder is a function that a composer (ex: Cond) will invoke
type Builder func()

// Cond constructs conditional branching
func Cond(e encoder.Type, cond, thenBranch, elseBranch Builder) {
	thenLabel := e.NewLabel()
	endLabel := e.NewLabel()

	cond()
	e.Append(isa.CondJump, thenLabel)
	elseBranch()
	e.Append(isa.Jump, endLabel)
	thenLabel.DropAnchor()
	thenBranch()
	endLabel.DropAnchor()
}
