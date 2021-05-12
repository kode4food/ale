package builtin_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/kode4food/ale/compiler/special"
	"github.com/kode4food/ale/core/bootstrap"
	"github.com/kode4food/ale/core/internal/builtin"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/read"
)

func interfaceErr(concrete, expected string) error {
	err := "interface conversion: %s is not %s: missing method"
	return fmt.Errorf(err, concrete, expected)
}

func TestApply(t *testing.T) {
	as := assert.New(t)

	as.True(builtin.IsApply.Call(builtin.Vector))
	as.False(builtin.IsApply.Call(S("55")))

	v1 := builtin.Vector.Call(S("4"), S("5"), S("6"))
	v2 := builtin.Apply.Call(builtin.Vector, S("1"), S("2"), S("3"), v1)
	v3 := builtin.Apply.Call(builtin.Vector, v1)

	as.String(`["4" "5" "6"]`, v1)
	as.String(`["1" "2" "3" "4" "5" "6"]`, v2)
	as.String(`["4" "5" "6"]`, v3)
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
	bootstrap.Into(e)

	as.False(builtin.IsSpecial.Call(builtin.Str))
	as.True(builtin.IsApply.Call(builtin.Str))

	i, ok := e.GetRoot().Resolve("if")
	as.True(ok && i.IsBound())
	as.True(builtin.IsSpecial.Call(i.Value()))
	as.False(builtin.IsApply.Call(i.Value()))
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
			(let [foo (lambda () greeting)]
				(call foo)))
	`, S("hello"))
}

func TestBadLambdaEval(t *testing.T) {
	as := assert.New(t)

	eNum := fmt.Errorf(special.ErrUnexpectedLambdaSyntax, "99")
	as.PanicWith(`(lambda-rec 99 "hello")`, eNum)

	eSym := fmt.Errorf(special.ErrUnexpectedLambdaSyntax, "foo/bar")
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

	e := interfaceErr("data.Integer", "data.Function")
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
	`, errors.New(read.ErrInvalidListSyntax))

	as.PanicWith(`
		(lambda (x y . z g) "explode")
	`, errors.New(read.ErrInvalidListSyntax))

	as.PanicWith(`
		(lambda (x y . . z) "explode")
	`, errors.New(read.ErrInvalidListSyntax))
}

func TestTailCallEval(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping tail call tests")
		return
	}
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
