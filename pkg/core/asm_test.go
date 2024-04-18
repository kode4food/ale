package core_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/pkg/compiler/encoder"
	"github.com/kode4food/ale/pkg/compiler/ir/analysis"
	builtin "github.com/kode4food/ale/pkg/core"
	"github.com/kode4food/ale/pkg/runtime/isa"
)

func TestAsmAddition(t *testing.T) {
	as := assert.New(t)
	as.EvalTo(`
		(define* test
			(lambda () (asm* pos-int 1 pos-int 2 add)))
		(test)
	`, I(3))
}

func TestAsmConstant(t *testing.T) {
	as := assert.New(t)
	as.EvalTo(
		`(asm* .const ("this is a list" 1 2 3))`,
		L(S("this is a list"), I(1), I(2), I(3)),
	)

	as.EvalTo(
		`(asm* .const 1 .const 2 .const 3 add add)`,
		I(6),
	)
}

func TestAsmJump(t *testing.T) {
	as := assert.New(t)
	as.EvalTo(`
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

func TestAsmLabelError(t *testing.T) {
	as := assert.New(t)

	as.PanicWith(`
		(asm*
			true
			cond-jump "not-a-label"
		:not-a-label)
    `, fmt.Errorf(builtin.ErrUnexpectedLabel, "not-a-label"))
}

func TestAsmLabelNumbering(t *testing.T) {
	as := assert.New(t)

	as.EncodesAs(isa.Instructions{
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

func TestAsmOutOfScopeError(t *testing.T) {
	as := assert.New(t)
	as.PanicWith(`
		(asm*
			.push-locals
			.local wont-be-found :val
			.const "hello"
			store wont-be-found
			.pop-locals
			load wont-be-found)
    `, fmt.Errorf(builtin.ErrUnexpectedName, "wont-be-found"))
}

func TestAsmLocalScopeError(t *testing.T) {
	as := assert.New(t)
	as.PanicWith(`
		(asm*
			.pop-locals
			.local hello :val)
	`, errors.New(encoder.ErrNoLocalScope))
}

func TestAsmValue(t *testing.T) {
	as := assert.New(t)
	as.EvalTo(`(asm* .eval (+ 1 2))`, I(3))
	as.EncodesAs(isa.Instructions{
		isa.PosInt.New(2),
		isa.PosInt.New(1),
		isa.Const.New(0),
		isa.Call.New(2),
		isa.Return.New(),
	}, `
	(asm*
		.eval (+ 1 2)
		return)
	`)
}

func TestAsmMakeEncoder(t *testing.T) {
	as := assert.New(t)

	as.EvalTo(`
		(define* if'
			(asm* !make-special
				[(predicate consequent alternative)
					.eval predicate
					cond-jump :consequent
					.eval alternative
					jump :end
				:consequent
					.eval consequent
				:end]
				[(predicate consequent)
					.eval predicate
					cond-jump :consequent
					nil
					jump :end
				:consequent
					.eval consequent
				:end]))

		(if' true "yep" "nope")
    `, S("yep"))
}

func TestAsmMakeRestEncoder(t *testing.T) {
	as := assert.New(t)

	as.EvalTo(`
		(define* test
			(asm* !make-special
				[(head . rest)
					.eval head]))
		[(test 1 2 3 4) (test 5 6) (test 7)]
	`, V(I(1), I(5), I(7)))

	as.EvalTo(`
		(define* test
			(asm* !make-special
				[(head . rest)
					.eval rest]))
		[(test 1 2 3 4) (test 5 6) (test 7)]
	`, V(V(I(2), I(3), I(4)), V(I(6)), V()))
}

func TestAsmOperandSizeError(t *testing.T) {
	as := assert.New(t)
	as.EvalTo(
		fmt.Sprintf("(asm* pos-int %d)", isa.OperandMask),
		I(int64(isa.OperandMask)),
	)

	as.PanicWith(
		fmt.Sprintf("(asm* pos-int %d)", isa.OperandMask+1),
		fmt.Errorf(isa.ErrExpectedOperand, isa.OperandMask+1),
	)
}

func TestAsmStackSizeError(t *testing.T) {
	as := assert.New(t)
	as.PanicWith(`(asm* pop)`, fmt.Errorf(analysis.ErrBadStackTermination, -2))
}
