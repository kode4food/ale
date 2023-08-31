package isa_test

import (
	"testing"

	"github.com/kode4food/ale/internal/assert"
	"github.com/kode4food/ale/runtime/isa"
)

func TestFlattenJump(t *testing.T) {
	as := assert.New(t)

	i1 := isa.Flatten(isa.Instructions{
		isa.NoOp.New(),
		isa.Label.New(0),
		isa.Jump.New(0),
	})

	as.Equal(isa.Instructions{
		isa.Jump.New(0),
	}, i1)
}

func TestFlattenCondJump(t *testing.T) {
	as := assert.New(t)

	i1 := isa.Flatten(isa.Instructions{
		isa.NoOp.New(),
		isa.Label.New(0),
		isa.NoOp.New(),
		isa.False.New(),
		isa.NoOp.New(),
		isa.CondJump.New(0),
		isa.NoOp.New(),
	})

	as.Equal(isa.Instructions{
		isa.False.New(),
		isa.CondJump.New(0),
	}, i1)
}

func TestForwardJump(t *testing.T) {
	as := assert.New(t)

	i1 := isa.Flatten(isa.Instructions{
		isa.NoOp.New(),
		isa.Label.New(0),
		isa.NoOp.New(),
		isa.Jump.New(1),
		isa.NoOp.New(),
		isa.Label.New(1),
		isa.NoOp.New(),
		isa.Jump.New(0),
	})

	as.Equal(isa.Instructions{
		isa.Jump.New(1),
		isa.Jump.New(0),
	}, i1)
}

func TestDoubleAnchor(t *testing.T) {
	as := assert.New(t)

	defer as.ExpectPanic(isa.ErrLabelAlreadyAnchored)

	isa.Flatten(isa.Instructions{
		isa.Label.New(0),
		isa.Label.New(0),
	})
}
