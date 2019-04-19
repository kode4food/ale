package test

import (
	"fmt"
	"testing"

	"gitlab.com/kode4food/ale/compiler/special"
	"gitlab.com/kode4food/ale/data"
	. "gitlab.com/kode4food/ale/internal/assert/helpers"
)

func TestFunctionPredicates(t *testing.T) {
	testCode(t, `(apply? if)`, data.False)
	testCode(t, `(!apply? if)`, data.True)
	testCode(t, `(special? def)`, data.True)
	testCode(t, `(!special? def)`, data.False)
	testCode(t, `(apply? 99)`, data.False)
	testCode(t, `(!apply? 99)`, data.True)
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
	e := typeErr("data.Integer", "*data.List")
	testBadCode(t, `(fn 99 "hello")`, e)

	e = interfaceErr("data.qualifiedSymbol", "data.LocalSymbol", "LocalSymbol")
	testBadCode(t, `(fn foo/bar [] "hello")`, e)
}

func TestApply(t *testing.T) {
	testCode(t, `(apply + [1 2 3])`, F(6))
	testCode(t, `
		(apply
			(fn add [x y z] (+ x y z))
			[1 2 3])
	`, F(6))

	e := interfaceErr("data.Integer", "data.Caller", "Caller")
	testBadCode(t, `(apply 32 [1 2 3])`, e)
}

func TestRestFunctions(t *testing.T) {
	testCode(t, `
		(def test (fn [f & r] (apply vector (cons f r))))
		(test 1 2 3 4 5 6 7)
	`, data.String("[1 2 3 4 5 6 7]"))

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
