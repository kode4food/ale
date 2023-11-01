package builtin_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/kode4food/ale/core/bootstrap"
	"github.com/kode4food/ale/core/builtin"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/read/parse"
	"github.com/kode4food/ale/runtime"
)

func unexpectedTypeError(got, expected string) error {
	return fmt.Errorf(runtime.ErrUnexpectedType, got, expected)
}

func TestApply(t *testing.T) {
	as := assert.New(t)

	as.True(getPredicate(builtin.FunctionKey).Call(builtin.Vector))
	as.False(getPredicate(builtin.FunctionKey).Call(S("55")))

	as.EvalTo(`
		(apply + '(1 2 3))`, I(6))

	as.EvalTo(`
		(apply + 9 8 7 '(1 2 3))`, I(30))
}

func TestPartialEval(t *testing.T) {
	as := assert.New(t)

	as.EvalTo(`
		(let [plus3 (partial +)]
			(plus3 1 1 1))`, I(3))

	as.EvalTo(`
		(let [plus3 (partial + 1 2)]
			(plus3 1 1 1))`, I(6))
}

func TestFunctionPredicates(t *testing.T) {
	as := assert.New(t)

	e := bootstrap.DevNullEnvironment()

	as.False(getPredicate(builtin.SpecialKey).Call(builtin.Str))
	as.True(getPredicate(builtin.FunctionKey).Call(builtin.Str))

	i, ok := e.GetRoot().Resolve("if")
	as.True(ok && i.IsBound())
	as.True(getPredicate(builtin.SpecialKey).Call(i.Value()))
	as.False(getPredicate(builtin.FunctionKey).Call(i.Value()))
}

func TestFunctionPredicatesEval(t *testing.T) {
	as := assert.New(t)
	as.EvalTo(`(apply? if)`, data.False)
	as.EvalTo(`(!apply? if)`, data.True)
	as.EvalTo(`(special? define*)`, data.True)
	as.EvalTo(`(!special? define*)`, data.False)
	as.EvalTo(`(apply? 99)`, data.False)
	as.EvalTo(`(!apply? 99)`, data.True)
}

func TestLambdaEval(t *testing.T) {
	as := assert.New(t)
	as.EvalTo(`
		(define call (lambda (func) (func)))
		(let [greeting "hello"]
			(let [foo (thunk greeting)]
				(call foo)))
	`, S("hello"))
}

func TestBadLambdaEval(t *testing.T) {
	as := assert.New(t)

	eNum := fmt.Errorf(builtin.ErrUnexpectedCaseSyntax, "99")
	as.PanicWith(`(lambda-rec 99 "hello")`, eNum)

	eSym := fmt.Errorf(builtin.ErrUnexpectedCaseSyntax, "foo/bar")
	as.PanicWith(`(lambda-rec foo/bar () "hello")`, eSym)
}

func TestApplyEval(t *testing.T) {
	as := assert.New(t)
	as.EvalTo(`(apply + [1 2 3])`, F(6))
	as.EvalTo(`
		(apply
			(lambda-rec add (x y z) (+ x y z))
			[1 2 3])
	`, F(6))

	e := unexpectedTypeError("integer", "lambda")
	as.PanicWith(`(apply 32 [1 2 3])`, e)
}

func TestRestFunctionsEval(t *testing.T) {
	as := assert.New(t)
	as.EvalTo(`
		(define test (lambda (f . r) (apply vector (cons f r))))
		(test 1 2 3 4 5 6 7)
	`, S("[1 2 3 4 5 6 7]"))

	as.PanicWith(`
		(lambda (x y .) "explode")
	`, errors.New(parse.ErrInvalidListSyntax))

	as.PanicWith(`
		(lambda (x y . z g) "explode")
	`, errors.New(parse.ErrInvalidListSyntax))

	as.PanicWith(`
		(lambda (x y . . z) "explode")
	`, errors.New(parse.ErrInvalidListSyntax))
}

func TestTailCallEval(t *testing.T) {
	as := assert.New(t)
	as.EvalTo(`
		(define-lambda to-zero (x)
			(cond
				[(> x 1000) (to-zero (- x 1))]
				[(> x 0)    (to-zero (- x 1))]
				[:else 0]))

		(to-zero 9999999)
	`, I(0))
}
