package encoder_test

import (
	"testing"

	"github.com/kode4food/ale/runtime/isa"
)

func TestLabels(t *testing.T) {
	as := NewWrapped(t)

	e := getTestEncoder()
	l1 := e.NewLabel()
	l2 := e.NewLabel()
	e.Emit(isa.Jump, l2)
	e.Emit(isa.NoOp)
	e.Emit(isa.Jump, l1)
	l2.DropAnchor()
	e.Emit(isa.NoOp)
	l1.DropAnchor()

	as.Instructions(isa.Instructions{
		&isa.Instruction{
			Opcode: isa.Jump,
			Args:   []isa.Word{1},
		},
		&isa.Instruction{
			Opcode: isa.NoOp,
		},
		&isa.Instruction{
			Opcode: isa.Jump,
			Args:   []isa.Word{0},
		},
		&isa.Instruction{
			Opcode: isa.Label,
			Args:   []isa.Word{1},
		},
		&isa.Instruction{
			Opcode: isa.NoOp,
		},
		&isa.Instruction{
			Opcode: isa.Label,
			Args:   []isa.Word{0},
		},
	}, e.Code())
}
