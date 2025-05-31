package special_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/internal/lang/params"
	"github.com/kode4food/ale/internal/lang/parse"
	"github.com/kode4food/ale/internal/runtime"
	"github.com/kode4food/ale/pkg/core/bootstrap"
	"github.com/kode4food/ale/pkg/core/builtin"
	"github.com/kode4food/ale/pkg/data"
)

func unexpectedTypeError(got, expected string) error {
	return fmt.Errorf(runtime.ErrUnexpectedType, got, expected)
}

func getPredicate(kwd data.Keyword) data.Procedure {
	return builtin.IsA.Call(kwd).(data.Procedure)
}

func TestApply(t *testing.T) {
	as := assert.New(t)

	as.True(getPredicate(builtin.ProcedureKey).Call(builtin.Vector))
	as.False(getPredicate(builtin.ProcedureKey).Call(S("55")))

	as.MustEvalTo(`
		(apply + '(1 2 3))`, I(6))

	as.MustEvalTo(`
		(apply + 9 8 7 '(1 2 3))`, I(30))
}

func TestPartialEval(t *testing.T) {
	as := assert.New(t)

	as.MustEvalTo(`
		(let [plus3 (partial +)]
			(plus3 1 1 1))`, I(3))

	as.MustEvalTo(`
		(let [plus3 (partial + 1 2)]
			(plus3 1 1 1))`, I(6))
}

func TestFunctionPredicates(t *testing.T) {
	as := assert.New(t)

	as.False(getPredicate(builtin.SpecialKey).Call(builtin.Str))
	as.True(getPredicate(builtin.ProcedureKey).Call(builtin.Str))

	e := bootstrap.DevNullEnvironment()
	root := e.GetRoot()
	as.True(getPredicate(builtin.SpecialKey).Call(as.IsBound(root, "if")))
	as.False(getPredicate(builtin.ProcedureKey).Call(as.IsBound(root, "if")))
}

func TestProcedurePredicatesEval(t *testing.T) {
	as := assert.New(t)
	as.MustEvalTo(`(procedure? if)`, data.False)
	as.MustEvalTo(`(!procedure? if)`, data.True)
	as.MustEvalTo(`(special? define*)`, data.True)
	as.MustEvalTo(`(!special? define*)`, data.False)
	as.MustEvalTo(`(procedure? 99)`, data.False)
	as.MustEvalTo(`(!procedure? 99)`, data.True)
}

func TestLambdaEval(t *testing.T) {
	as := assert.New(t)
	as.MustEvalTo(`
		(define call (lambda (func) (func)))
		(let [greeting "hello"]
			(let [foo (thunk greeting)]
				(call foo)))
	`, S("hello"))
}

func TestBadLambdaEval(t *testing.T) {
	as := assert.New(t)

	eNum := fmt.Errorf(params.ErrUnexpectedCaseSyntax, "99")
	as.ErrorWith(`(lambda-rec 99 "hello")`, eNum)

	eSym := fmt.Errorf(params.ErrUnexpectedCaseSyntax, "foo/bar")
	as.ErrorWith(`(lambda-rec foo/bar () "hello")`, eSym)
}

func TestApplyEval(t *testing.T) {
	as := assert.New(t)
	as.MustEvalTo(`(apply + [1 2 3])`, F(6))
	as.MustEvalTo(`
		(apply
			(lambda-rec add (x y z) (+ x y z))
			[1 2 3])
	`, F(6))

	e := unexpectedTypeError("integer", "procedure")
	as.PanicWith(`(apply 32 [1 2 3])`, e)
}

func TestRestFunctionsEval(t *testing.T) {
	as := assert.New(t)
	as.MustEvalTo(`
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
	as.MustEvalTo(`
		(define-lambda to-zero (x)
			(cond
				[(> x 1000) (to-zero (- x 1))]
				[(> x 0)    (to-zero (- x 1))]
				[:else 0]))

		(to-zero 999999)
	`, I(0))
}
