package special_test

import (
	"fmt"
	"testing"

	"github.com/kode4food/ale/compiler/encoder"
	"github.com/kode4food/ale/compiler/special"
	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/runtime/isa"
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
	defer as.ExpectPanic(
		fmt.Sprintf(special.ErrUnexpectedLabel, "not-a-label"),
	)
	as.Eval(`
		(asm*
			true
			cond-jump not-a-label
		:not-a-label)
    `)
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
	defer as.ExpectPanic(
		fmt.Sprintf(special.ErrUnexpectedName, "wont-be-found"),
	)
	as.Eval(`
		(asm*
			.push-locals
			.local wont-be-found :val
			.const "hello"
			store wont-be-found
			.pop-locals
			load wont-be-found)
    `)
}

func TestAsmLocalScopeError(t *testing.T) {
	as := assert.New(t)
	defer as.ExpectPanic(encoder.ErrNoLocalScope)
	as.Eval(`
		(asm*
			.pop-locals
			.local hello :val)
	`)
}

func TestAsmValue(t *testing.T) {
	as := assert.New(t)
	as.EvalTo(`(asm* .value (+ 1 2))`, I(3))
	as.EncodesAs(isa.Instructions{
		isa.PosInt.New(2),
		isa.PosInt.New(1),
		isa.Const.New(0),
		isa.Call.New(2),
		isa.Return.New(),
	}, `
	(asm*
		.value (+ 1 2)
		return)
	`)
}

func TestAsmMakeEncoder(t *testing.T) {
	as := assert.New(t)

	as.EvalTo(`
		(define* if'
			(asm* !make-encoder
				[(predicate consequent alternative)
					.value predicate
					make-truthy
					cond-jump :consequent
					.value alternative
					jump :end
				:consequent
					.value consequent
				:end]
				[(predicate consequent)
					.value predicate
					make-truthy
					cond-jump :consequent
					nil
					jump :end
				:consequent
					.value consequent
				:end]))

		(if' true "yep" "nope")
    `, S("yep"))
}

func TestAsmMakeRestEncoder(t *testing.T) {
	as := assert.New(t)

	as.EvalTo(`
		(define* test
			(asm* !make-encoder
				[(head . rest)
					.value head]))
		(test 1 2 3 4)
	`, I(1))

	as.EvalTo(`
		(define* test
			(asm* !make-encoder
				[(head . rest)
					.value rest]))
		(test 1 2 3 4)
	`, V(I(2), I(3), I(4)))
}

func TestAsmOperandSizeError(t *testing.T) {
	as := assert.New(t)
	as.EvalTo(
		fmt.Sprintf("(asm* pos-int %d)", isa.OperandMask),
		I(isa.OperandMask),
	)

	defer as.ExpectPanic(
		fmt.Sprintf(isa.ErrExpectedOperand, isa.OperandMask+1),
	)
	as.Eval(fmt.Sprintf("(asm* pos-int %d)", isa.OperandMask+1))
}
