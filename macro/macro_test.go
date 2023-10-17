package macro_test

import (
	"testing"

	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/macro"
)

func TestMacroCall(t *testing.T) {
	as := assert.New(t)

	d, ok := as.Eval(`define`).(macro.Call)
	as.True(ok)
	as.NotNil(d)

	as.Equal(macro.CallType, d.Type())
	as.False(d.Equal(d))
	as.Contains(":type macro", d)
}

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
