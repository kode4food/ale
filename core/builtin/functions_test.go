package builtin_test

import (
	"fmt"
	"testing"

	"github.com/kode4food/ale/compiler/special"
	"github.com/kode4food/ale/core/bootstrap"
	"github.com/kode4food/ale/core/builtin"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/read"
)

func interfaceErr(concrete, intf, method string) error {
	err := "interface conversion: %s is not %s: missing method %s"
	return fmt.Errorf(err, concrete, intf, method)
}

func typeErr(concrete, expected string) error {
	err := "interface conversion: data.Value is %s, not %s"
	return fmt.Errorf(err, concrete, expected)
}

func getCall(v data.Value) data.Call {
	return v.(data.Caller).Caller()
}

func TestApply(t *testing.T) {
	as := assert.New(t)

	vCall := data.MakeApplicative(builtin.Vector, nil)
	as.True(builtin.IsApply(vCall))
	as.False(builtin.IsApply(S("55")))

	v1 := builtin.Vector(S("4"), S("5"), S("6"))
	v2 := builtin.Apply(vCall, S("1"), S("2"), S("3"), v1)
	v3 := builtin.Apply(vCall, v1)

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

	manager := bootstrap.DevNullManager()
	bootstrap.Into(manager)

	f1 := data.MakeApplicative(builtin.Str, nil)
	as.False(builtin.IsSpecial(f1))
	as.True(builtin.IsApply(f1))

	e, ok := manager.GetRoot().Resolve("if")
	as.True(ok && e.IsBound())
	as.True(builtin.IsSpecial(e.Value()))
	as.False(builtin.IsApply(e.Value()))
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

	eNum := fmt.Errorf(special.UnexpectedLambdaSyntax, "99")
	as.PanicWith(`(lambda-rec 99 "hello")`, eNum)

	eSym := fmt.Errorf(special.UnexpectedLambdaSyntax, "foo/bar")
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

	e := interfaceErr("data.Integer", "data.Caller", "Caller")
	as.PanicWith(`(apply 32 [1 2 3])`, e)
}

func TestRestFunctionsEval(t *testing.T) {
	as := assert.New(t)
	as.EvalTo(`
		(define test (lambda (f . r) (apply vector (cons f r))))
		(test 1 2 3 4 5 6 7)
	`, data.String("[1 2 3 4 5 6 7]"))

	as.PanicWith(`
		(lambda (x y .) "explode")
	`, fmt.Errorf(read.InvalidListSyntax))

	as.PanicWith(`
		(lambda (x y . z g) "explode")
	`, fmt.Errorf(read.InvalidListSyntax))

	as.PanicWith(`
		(lambda (x y . . z) "explode")
	`, fmt.Errorf(read.InvalidListSyntax))
}

func TestTailCallEval(t *testing.T) {
	as := assert.New(t)
	as.EvalTo(`
		(define-lambda to-zero (x)
			(cond
				[(> x 1000) (to-zero (- x 1))]
				[(> x 0)    (to-zero (- x 1))]
				[:else 0]))

		(to-zero 9999)
	`, data.Integer(0))
}
