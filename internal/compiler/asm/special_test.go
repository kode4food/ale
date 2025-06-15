package asm_test

import (
	"testing"

	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
)

func TestMakeSpecial(t *testing.T) {
	as := assert.New(t)

	as.MustEvalTo(`
		(define* if'
			(special*
				[(predicate consequent alternative)
					eval predicate
					cond-jump :consequent
					eval alternative
					jump :end
				:consequent
					eval consequent
				:end]
				[(predicate consequent)
					eval predicate
					cond-jump :consequent
					null
					jump :end
				:consequent
					eval consequent
				:end]))

		(if' true "yep" "nope")
    `, S("yep"))
}

func TestMakeRestSpecial(t *testing.T) {
	as := assert.New(t)

	as.MustEvalTo(`
        (define* test
			(special*
				[(head . rest)
					eval head]))
		[(test 1 2 3 4) (test 5 6) (test 7)]
	`, V(I(1), I(5), I(7)))

	as.MustEvalTo(`
		(define* test
			(special*
				[(head . rest)
					eval rest]))
		[(test 1 2 3 4) (test 5 6) (test 7)]
	`, V(V(I(2), I(3), I(4)), V(I(6)), V()))
}
