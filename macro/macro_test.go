package macro_test

import (
	"testing"

	"github.com/kode4food/ale/internal/assert"
	"github.com/kode4food/ale/macro"
	"github.com/kode4food/ale/read"
)

func TestMacroCall(t *testing.T) {
	as := assert.New(t)

	d, ok := as.MustEval(`define`).(macro.Call)
	as.True(ok)
	if as.NotNil(d) {
		as.False(d.Type().Accepts(macro.CallType))
		as.True(macro.CallType.Accepts(d.Type()))
		as.False(d.Equal(d))
		as.Contains(":type macro", d)
	}
}

func TestExpand(t *testing.T) {
	as := assert.New(t)

	ns := assert.GetTestNamespace()
	as.MustEvalTo(
		`(macroexpand '(define (name . _) "hello"))`,
		read.MustFromString(ns,
			`(ale/%define name (ale/label name (ale/lambda _ "hello")))`,
		).Car(),
	)
}

func TestExpand1(t *testing.T) {
	as := assert.New(t)

	ns := assert.GetTestNamespace()
	as.MustEvalTo(
		`(macroexpand-1 '(define (name . _) (or false true)))`,
		read.MustFromString(ns,
			`(ale/%define name (ale/label name (ale/lambda _ (or false true))))`,
		).Car(),
	)
}
