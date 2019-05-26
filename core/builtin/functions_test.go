package builtin_test

import (
	"fmt"
	"testing"

	"gitlab.com/kode4food/ale/compiler/special"
	"gitlab.com/kode4food/ale/core/bootstrap"
	"gitlab.com/kode4food/ale/core/builtin"
	"gitlab.com/kode4food/ale/data"
	"gitlab.com/kode4food/ale/internal/assert"
	. "gitlab.com/kode4food/ale/internal/assert/helpers"
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

	vCall := data.Call(builtin.Vector)
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

	f1 := data.ApplicativeFunction(builtin.Str)
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
	as.EvalTo(`(special? def)`, data.True)
	as.EvalTo(`(!special? def)`, data.False)
	as.EvalTo(`(apply? 99)`, data.False)
	as.EvalTo(`(!apply? 99)`, data.True)
}

func TestLambdaEval(t *testing.T) {
	as := assert.New(t)
	as.EvalTo(`
		(def call (fn [func] (func)))
		(let [greeting "hello"]
			(let [foo (fn [] greeting)]
				(call foo)))
	`, S("hello"))
}

func TestBadLambdaEval(t *testing.T) {
	as := assert.New(t)

	e := typeErr("data.Integer", "*data.List")
	as.PanicWith(`(fn 99 "hello")`, e)

	e = interfaceErr("data.qualifiedSymbol", "data.LocalSymbol", "LocalSymbol")
	as.PanicWith(`(fn foo/bar [] "hello")`, e)
}

func TestApplyEval(t *testing.T) {
	as := assert.New(t)
	as.EvalTo(`(apply + [1 2 3])`, F(6))
	as.EvalTo(`
		(apply
			(fn add [x y z] (+ x y z))
			[1 2 3])
	`, F(6))

	e := interfaceErr("data.Integer", "data.Caller", "Caller")
	as.PanicWith(`(apply 32 [1 2 3])`, e)
}

func TestRestFunctionsEval(t *testing.T) {
	as := assert.New(t)
	as.EvalTo(`
		(def test (fn [f & r] (apply vector (cons f r))))
		(test 1 2 3 4 5 6 7)
	`, data.String("[1 2 3 4 5 6 7]"))

	as.PanicWith(`
		(fn [x y &] "explode")
	`, fmt.Errorf(special.InvalidRestArgument, "[]"))

	as.PanicWith(`
		(fn [x y & z g] "explode")
	`, fmt.Errorf(special.InvalidRestArgument, "[z g]"))

	as.PanicWith(`
		(fn [x y & & z] "explode")
	`, fmt.Errorf(special.InvalidRestArgument, "[& z]"))
}

func TestTailCallEval(t *testing.T) {
	as := assert.New(t)
	as.EvalTo(`
		(defn to-zero [x]
			(cond
				(> x 1000) (to-zero (- x 1))
				(> x 0)    (to-zero (- x 1))
				:else 0))

		(to-zero 9999)
	`, data.Integer(0))
}
