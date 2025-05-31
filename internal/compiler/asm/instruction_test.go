package asm_test

import (
	"fmt"
	"testing"

	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/internal/compiler/asm"
	"github.com/kode4food/ale/internal/runtime/isa"
)

func TestAddition(t *testing.T) {
	as := assert.New(t)
	as.MustEvalTo(`
		(define* test
			(lambda () (asm* pos-int 1 pos-int 2 add)))
		(test)
	`, I(3))
}

func TestJump(t *testing.T) {
	as := assert.New(t)
	as.MustEvalTo(`
		(define* test
			(lambda ()
				(asm*
					.local some-value :val
					true
					store some-value
					load some-value
					cond-jump :first
					pos-int 0
					jump :second
				:first
					pos-int 1
				:second)))
		(test)
    `, I(1))
}

func TestLabelError(t *testing.T) {
	as := assert.New(t)

	as.ErrorWith(`
		(asm*
			true
			cond-jump "not-a-label"
		:not-a-label)
    `, fmt.Errorf(asm.ErrUnexpectedLabel, "not-a-label"))
}

func TestLabelNumbering(t *testing.T) {
	as := assert.New(t)

	as.MustEncodedAs(isa.Instructions{
		isa.Jump.New(0),
		isa.NoOp.New(),
		isa.Jump.New(1),
		isa.Label.New(0),
		isa.NoOp.New(),
		isa.Label.New(1),
	}, `
		(asm*
			jump :second
			no-op
			jump :first
		:second
			no-op
		:first)
	`)
}

func TestOperandSizeError(t *testing.T) {
	as := assert.New(t)
	as.MustEvalTo(
		fmt.Sprintf("(asm* pos-int %d)", isa.OperandMask),
		I(int64(isa.OperandMask)),
	)

	as.ErrorWith(
		fmt.Sprintf("(asm* pos-int %d)", isa.OperandMask+1),
		fmt.Errorf(isa.ErrExpectedOperand, isa.OperandMask+1),
	)
}
