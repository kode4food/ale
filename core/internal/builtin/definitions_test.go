package builtin_test

import (
	"fmt"
	"testing"

	"github.com/kode4food/ale/compiler/special"
	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
)

func TestFunctionEval(t *testing.T) {
	as := assert.New(t)
	as.EvalTo(`
		(define-lambda say-hello ()
		  "Hello, World!")
		(say-hello)
	`, S("Hello, World!"))

	as.EvalTo(`
		(define-lambda identity (value) value)
		(identity "foo")
	`, S("foo"))
}

func TestBadFunctionEval(t *testing.T) {
	symErr := interfaceErr("data.Integer", "data.LocalSymbol")
	numErr := fmt.Errorf(special.ErrUnexpectedLambdaSyntax, "99")
	vecErr := interfaceErr("data.Integer", "data.Vector")
	listErr := interfaceErr("data.Integer", "data.LocalSymbol")

	as := assert.New(t)
	as.PanicWith(`(define-lambda blah (name 99 bad) (name))`, symErr)
	as.PanicWith(`(define-lambda blah 99 (name))`, numErr)
	as.PanicWith(`(define-lambda 99 (x y) (+ x y))`, symErr)
	as.PanicWith(`(define-lambda blah (99 "hello"))`, listErr)
	as.PanicWith(`(define-lambda blah [(x) "hello"] 99)`, vecErr)
}
