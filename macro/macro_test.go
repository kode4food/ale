package macro_test

import (
	"testing"

	"github.com/kode4food/ale/internal/assert"
	"github.com/kode4food/ale/macro"
	"github.com/kode4food/ale/read"
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
		read.FromString(
			`(ale/define* name (ale/label name (ale/lambda _ "hello")))`,
		).Car(),
	)
}

func TestExpand1(t *testing.T) {
	as := assert.New(t)

	as.EvalTo(
		`(macroexpand-1 '(define (name . _) (or false true)))`,
		read.FromString(
			`(ale/define* name (ale/label name (ale/lambda _ (or false true))))`,
		).Car(),
	)
}
