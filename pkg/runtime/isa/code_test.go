package isa_test

import (
	"fmt"
	"testing"

	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/pkg/runtime/isa"
)

func TestInstructions(t *testing.T) {
	as := assert.New(t)

	i1 := isa.CondJump.New(27)
	as.String("cond-jump 27", i1)

	defer func() {
		if rec := recover(); rec != nil {
			err := fmt.Errorf(isa.ErrBadInstruction, isa.CondJump)
			as.Equal(err, rec)
		} else {
			as.Fail("proper error not raised")
		}
	}()

	isa.CondJump.New(12, 32)
}

func TestInstructionsString(t *testing.T) {
	as := assert.New(t)
	inst := isa.Instructions{
		isa.Const.New(0),
		isa.CondJump.New(2),
		isa.Return.New(),
	}
	as.String("const 0\ncond-jump 2\nreturn", inst.String())
}

func TestInstructionEquality(t *testing.T) {
	as := assert.New(t)
	i1 := isa.Const.New(0)
	i2 := isa.Const.New(0) // Some content
	i3 := isa.Const.New(1) // Different Arg
	i4 := isa.Load.New(0)  // Different Opcode

	as.True(i1.Equal(i1))
	as.True(i1.Equal(i2))
	as.False(i1.Equal(i3))
	as.False(i1.Equal(i4))
	as.False(i1.Equal(I(37)))
}

func TestInstructionSplit(t *testing.T) {
	as := assert.New(t)
	cj := isa.CondJump.New(37)
	oc, op := cj.Split()
	as.Equal(isa.CondJump, oc)
	as.Equal(isa.Operand(37), op)
}

func TestInstructionOperandSizeError(t *testing.T) {
	as := assert.New(t)
	defer as.ExpectPanic(
		fmt.Errorf(isa.ErrExpectedOperand, isa.OperandMask+1),
	)
	isa.PosInt.New(isa.OperandMask + 1)
}
