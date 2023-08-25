package special_test

import (
	"fmt"
	"testing"

	"github.com/kode4food/ale/compiler/special"

	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
)

func TestAsmAddition(t *testing.T) {
	as := assert.New(t)
	as.EvalTo(`
		(define* test
			(lambda () (asm* one two add)))
		(test)
	`, I(3))
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

	defer as.ExpectPanic(
		fmt.Sprintf(special.ErrUnexpectedLabel, "not-a-label"),
	)
	as.EvalTo(`
		(define* test
			(lambda () (asm*
				true
				cond-jump not-a-label
			:not-a-label)))
		(test)
    `, I(1))
}
