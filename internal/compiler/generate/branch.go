package generate

import (
	"github.com/kode4food/ale/internal/compiler/encoder"
	"github.com/kode4food/ale/internal/runtime/isa"
)

// Branch constructs predicate, consequent, alternative branching
func Branch(
	e encoder.Encoder, predicate, consequent, alternative Builder,
) error {
	thenLabel := e.NewLabel()
	endLabel := e.NewLabel()

	if err := predicate(e); err != nil {
		return err
	}
	e.Emit(isa.CondJump, thenLabel)
	if err := alternative(e); err != nil {
		return err
	}
	e.Emit(isa.Jump, endLabel)
	e.Emit(isa.Label, thenLabel)
	if err := consequent(e); err != nil {
		return err
	}
	e.Emit(isa.Label, endLabel)
	return nil
}
