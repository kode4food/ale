package test

import (
	"fmt"
	"testing"

	"gitlab.com/kode4food/ale/api"
	. "gitlab.com/kode4food/ale/internal/assert/helpers"
	"gitlab.com/kode4food/ale/internal/compiler/special"
)

func TestFunctionPredicates(t *testing.T) {
	testCode(t, `(apply? if)`, api.False)
	testCode(t, `(!apply? if)`, api.True)
	testCode(t, `(special? def)`, api.True)
	testCode(t, `(!special? def)`, api.False)
	testCode(t, `(apply? 99)`, api.False)
	testCode(t, `(!apply? 99)`, api.True)
}

func TestLambda(t *testing.T) {
	testCode(t, `
		(def call (fn [func] (func)))
		(let [greeting "hello"]
			(let [foo (fn [] greeting)]
				(call foo)))
	`, S("hello"))
}

func TestBadLambda(t *testing.T) {
	e := typeErr("api.Integer", "*api.List")
	testBadCode(t, `(fn 99 "hello")`, e)

	e = interfaceErr("api.qualifiedSymbol", "api.LocalSymbol", "LocalSymbol")
	testBadCode(t, `(fn foo/bar [] "hello")`, e)
}

func TestApply(t *testing.T) {
	testCode(t, `(apply + [1 2 3])`, F(6))
	testCode(t, `
		(apply
			(fn add [x y z] (+ x y z))
			[1 2 3])
	`, F(6))

	e := interfaceErr("api.Integer", "api.Caller", "Caller")
	testBadCode(t, `(apply 32 [1 2 3])`, e)
}

func TestRestFunctions(t *testing.T) {
	testCode(t, `
		(def test (fn [f & r] (apply vector (cons f r))))
		(test 1 2 3 4 5 6 7)
	`, api.String("[1 2 3 4 5 6 7]"))

	testBadCode(t, `
		(fn [x y &] "explode")
	`, fmt.Errorf(special.InvalidRestArgument, "[]"))

	testBadCode(t, `
		(fn [x y & z g] "explode")
	`, fmt.Errorf(special.InvalidRestArgument, "[z g]"))

	testBadCode(t, `
		(fn [x y & & z] "explode")
	`, fmt.Errorf(special.InvalidRestArgument, "[& z]"))
}
