package builtin_test

import (
	"testing"

	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
)

func TestPairsEval(t *testing.T) {
	as := assert.New(t)

	as.EvalTo(
		`(cons 7 (cons 8 (cons (cons 10 11) 9)))`,
		S(`(7 8 (10 . 11) . 9)`),
	)

	as.EvalTo(
		`(car (cons 7 (cons 8 (cons (cons 10 11) 9))))`,
		I(7),
	)

	as.EvalTo(
		`(cdr (cons 7 (cons 8 (cons (cons 10 11) 9))))`,
		S(`(8 (10 . 11) . 9)`),
	)

	as.EvalTo(`(pair? (cons 7 8))`, B(true))
	as.EvalTo(`(pair? 99)`, B(false))
}
