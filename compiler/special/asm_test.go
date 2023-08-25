package special_test

import (
	"fmt"
	"testing"

	"github.com/kode4food/ale/compiler/special"
	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/runtime/isa"
)

func TestAsmAddition(t *testing.T) {
	as := assert.New(t)
	as.EvalTo(`
		(define* test
			(lambda () (asm* one two add)))
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
			(lambda () (asm*
				.local some-value :val
				true
				store some-value
				load some-value
				cond-jump :first
				zero
				jump :second
			:first
				one
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
		isa.New(isa.Jump, 0),
		isa.New(isa.NoOp),
		isa.New(isa.Jump, 1),
		isa.New(isa.Label, 0),
		isa.New(isa.NoOp),
		isa.New(isa.Label, 1),
	}, `(asm*
		jump :second
		no-op
		jump :first
 	:second
		no-op
	:first)`)
}

func TestPushLocalsError(t *testing.T) {
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
