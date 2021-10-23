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

	as.Equal([]isa.Word{
		isa.Word(isa.Jump), 0,
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

	as.Equal([]isa.Word{
		isa.Word(isa.False),
		isa.Word(isa.CondJump), 0,
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

	as.Equal([]isa.Word{
		isa.Word(isa.Jump), 2,
		isa.Word(isa.Jump), 0,
	}, i1)
}

func TestDoubleAnchor(t *testing.T) {
	as := assert.New(t)

	defer as.ExpectProgrammerError("label has already been anchored")

	isa.Flatten(isa.Instructions{
		isa.New(isa.Label, 0),
		isa.New(isa.Label, 0),
	})
}
