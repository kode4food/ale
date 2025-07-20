package builtin_test

import (
	"testing"

	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
)

func TestPairsEval(t *testing.T) {
	as := assert.New(t)

	as.MustEvalTo(
		`(cons 7 (cons 8 (cons (cons 10 11) 9)))`,
		S(`(7 8 (10 . 11) . 9)`),
	)

	as.MustEvalTo(
		`(car (cons 7 (cons 8 (cons (cons 10 11) 9))))`,
		I(7),
	)

	as.MustEvalTo(
		`(cdr (cons 7 (cons 8 (cons (cons 10 11) 9))))`,
		S(`(8 (10 . 11) . 9)`),
	)

	as.MustEvalTo(`(pair? (cons 7 8))`, B(true))
	as.MustEvalTo(`(pair? 99)`, B(false))
}
