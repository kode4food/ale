package macro_test

import (
	"testing"

	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
)

func TestExpand(t *testing.T) {
	as := assert.New(t)

	as.EvalTo(
		`(macroexpand '(define (name . _) "hello"))`,
		L(
			QS("ale", "define*"),
			LS("name"),
			L(
				QS("ale", "label"),
				LS("name"),
				L(
					QS("ale", "lambda"),
					LS("_"),
					S("hello"),
				),
			),
		),
	)
}
