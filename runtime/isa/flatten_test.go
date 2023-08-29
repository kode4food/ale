package isa_test

import (
	"testing"

	"github.com/kode4food/ale/internal/assert"
	"github.com/kode4food/ale/runtime/isa"
)

func TestFlattenJump(t *testing.T) {
	as := assert.New(t)

	i1 := isa.Flatten(isa.Instructions{
		isa.New(isa.NoOp),
		isa.New(isa.Label, 0),
		isa.New(isa.Jump, 0),
	})

	as.Equal(isa.Instructions{
		isa.Instruction(isa.Jump),
	}, i1)
}

func TestFlattenCondJump(t *testing.T) {
	as := assert.New(t)

	i1 := isa.Flatten(isa.Instructions{
		isa.New(isa.NoOp),
		isa.New(isa.Label, 0),
		isa.New(isa.NoOp),
		isa.New(isa.False),
		isa.New(isa.NoOp),
		isa.New(isa.CondJump, 0),
		isa.New(isa.NoOp),
	})

	as.Equal(isa.Instructions{
		isa.Instruction(isa.False),
		isa.Instruction(isa.CondJump),
	}, i1)
}

func TestForwardJump(t *testing.T) {
	as := assert.New(t)

	i1 := isa.Flatten(isa.Instructions{
		isa.New(isa.NoOp),
		isa.New(isa.Label, 0),
		isa.New(isa.NoOp),
		isa.New(isa.Jump, 1),
		isa.New(isa.NoOp),
		isa.New(isa.Label, 1),
		isa.New(isa.NoOp),
		isa.New(isa.Jump, 0),
	})

	as.Equal(isa.Instructions{
		isa.New(isa.Jump, 1),
		isa.New(isa.Jump, 0),
	}, i1)
}

func TestDoubleAnchor(t *testing.T) {
	as := assert.New(t)

	defer as.ExpectPanic(isa.ErrLabelAlreadyAnchored)

	isa.Flatten(isa.Instructions{
		isa.New(isa.Label, 0),
		isa.New(isa.Label, 0),
	})
}
