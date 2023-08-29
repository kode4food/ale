package isa_test

import (
	"fmt"
	"testing"

	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/runtime/isa"
)

func TestInstructions(t *testing.T) {
	as := assert.New(t)

	i1 := isa.New(isa.CondJump, 27)
	as.String("CondJump(27)", i1)

	defer func() {
		if rec := recover(); rec != nil {
			err := fmt.Errorf(isa.ErrBadInstruction, isa.CondJump)
			as.Equal(err, rec)
		} else {
			as.Fail("proper error not raised")
		}
	}()

	isa.New(isa.CondJump, 12, 32)
}

func TestInstructionString(t *testing.T) {
	as := assert.New(t)
	inst := isa.New(isa.Const, 0)
	as.String(`Const(0)`, inst.String())
}

func TestInstructionEquality(t *testing.T) {
	as := assert.New(t)
	i1 := isa.New(isa.Const, 0)
	i2 := isa.New(isa.Const, 0) // Some content
	i3 := isa.New(isa.Const, 1) // Different Arg
	i4 := isa.New(isa.Load, 0)  // Different Opcode

	as.True(i1.Equal(i1))
	as.True(i1.Equal(i2))
	as.False(i1.Equal(i3))
	as.False(i1.Equal(i4))
	as.False(i1.Equal(I(37)))
}

func TestInstructionSplit(t *testing.T) {
	as := assert.New(t)
	cj := isa.New(isa.CondJump, 37)
	oc, op := cj.Split()
	as.Equal(isa.CondJump, oc)
	as.Equal(isa.Operand(37), op)
}
