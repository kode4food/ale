package special_test

import (
	"fmt"
	"testing"

	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/internal/lang/params"
	"github.com/kode4food/ale/internal/runtime"
)

func TestFunctionEval(t *testing.T) {
	as := assert.New(t)
	as.MustEvalTo(`
		(define-lambda say-hello ()
		  "Hello, World!")
		(say-hello)
	`, S("Hello, World!"))

	as.MustEvalTo(`
		(define-lambda identity (value) value)
		(identity "foo")
	`, S("foo"))
}

func TestBadFunctionEval(t *testing.T) {
	symErr := fmt.Errorf(runtime.ErrUnexpectedType, "integer", "local")
	numErr := fmt.Errorf(params.ErrUnexpectedCaseSyntax, "99")
	vecErr := fmt.Errorf(runtime.ErrUnexpectedType, "integer", "vector")
	invalidErr := fmt.Errorf("got number, expected local: 99")

	as := assert.New(t)
	as.PanicWith(`(define-lambda blah (name 99 bad) (name))`, symErr)
	as.PanicWith(`(define-lambda blah 99 (name))`, numErr)
	as.PanicWith(`(define-lambda 99 (x y) (+ x y))`, invalidErr)
	as.PanicWith(`(define-lambda blah (99 "hello"))`, symErr)
	as.PanicWith(`(define-lambda blah [(x) "hello"] 99)`, vecErr)
}
