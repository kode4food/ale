package macro_test

import (
	"testing"

	"github.com/kode4food/ale/internal/assert"
	"github.com/kode4food/ale/pkg/macro"
	"github.com/kode4food/ale/pkg/read"
)

func TestMacroCall(t *testing.T) {
	as := assert.New(t)

	d, ok := as.MustEval(`define`).(macro.Call)
	as.True(ok)
	as.NotNil(d)

	as.Equal(macro.CallType, d.Type())
	as.False(d.Equal(d))
	as.Contains(":type macro", d)
}

func TestExpand(t *testing.T) {
	as := assert.New(t)

	as.MustEvalTo(
		`(macroexpand '(define (name . _) "hello"))`,
		read.FromString(
			`(ale/define* name (ale/label name (ale/lambda _ "hello")))`,
		).Car(),
	)
}

func TestExpand1(t *testing.T) {
	as := assert.New(t)

	as.MustEvalTo(
		`(macroexpand-1 '(define (name . _) (or false true)))`,
		read.FromString(
			`(ale/define* name (ale/label name (ale/lambda _ (or false true))))`,
		).Car(),
	)
}
