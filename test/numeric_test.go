package test

import (
	"fmt"
	"testing"

	"gitlab.com/kode4food/ale/data"
	. "gitlab.com/kode4food/ale/internal/assert/helpers"
)

func TestBasicNumber(t *testing.T) {
	testCode(t, `(+)`, I(0))
	testCode(t, `(*)`, I(1))
	testCode(t, `(+ 1 1)`, I(2))
	testCode(t, `(* 4 4)`, I(16))
	testCode(t, `(+ 5 4)`, I(9))
	testCode(t, `(* 12 3)`, I(36))
	testCode(t, `(- 10 4)`, I(6))
	testCode(t, `(- 10 4 2)`, I(4))
	testCode(t, `(- 5)`, I(-5))
	testCode(t, `(/ 10 2)`, I(5))
	testCode(t, `(/ 10 2 5)`, I(1))
	testCode(t, `(mod 10 3)`, I(1))
	testCode(t, `(mod 100 8 7)`, I(4))
}

func TestNestedNumber(t *testing.T) {
	testCode(t, `(/ 10 (- 5 3))`, I(5))
	testCode(t, `(* 5 (- 5 3))`, I(10))
	testCode(t, `(/ 10 (/ 6 3))`, I(5))
}

func TestNonNumber(t *testing.T) {
	testCode(t, `(is-pos-inf (/ 99.0 0))`, data.True)
	testCode(t, `(is-pos-inf 99)`, data.False)
	testCode(t, `(is-pos-inf "hello")`, data.False)

	testCode(t, `(is-neg-inf (/ -99.0 0))`, data.True)
	testCode(t, `(is-neg-inf -99)`, data.False)
	testCode(t, `(is-neg-inf "hello")`, data.False)

	testCode(t, `(is-nan 99)`, data.False)
	testCode(t, `(is-nan "hello")`, data.False)
}

func TestBadMaths(t *testing.T) {
	e := interfaceErr("data.String", "data.Number", "Add")

	testBadCode(t, `(+ 99 "hello")`, e)
	testBadCode(t, `(+ "hello")`, e)
}

func TestBadNumbers(t *testing.T) {
	testBadNumber := func(err string, ns string) {
		testBadCode(t, ns, fmt.Errorf(err, data.String(ns)))
	}

	testBadNumber(data.ExpectedInteger, "0xfkk")
	testBadNumber(data.ExpectedInteger, "0b01109")
	testBadNumber(data.ExpectedInteger, "123j-k")
	testBadNumber(data.ExpectedFloat, "1.2j-k")
	//testBadNumber(data.ExpectedRatio, "1/2p")
}

func TestCompare(t *testing.T) {
	testCode(t, `(= 1 1)`, data.True)
	testCode(t, `(= 1 1 1 1 '1 1 1)`, data.True)
	testCode(t, `(= 1 2)`, data.False)
	testCode(t, `(= 1 1 1 1 2 1 1 1)`, data.False)

	testCode(t, `(!= 1 1)`, data.False)
	testCode(t, `(!= 1 1 1 1 '1 1 1)`, data.False)
	testCode(t, `(!= 1 2)`, data.True)
	testCode(t, `(!= 1 1 1 1 2 1 1 1)`, data.True)

	testCode(t, `(> 1 1)`, data.False)
	testCode(t, `(> 2 1)`, data.True)
	testCode(t, `(> 1 2)`, data.False)
	testCode(t, `(> 1 2 3 4 5)`, data.False)
	testCode(t, `(> 5 4 3 2 1)`, data.True)
	testCode(t, `(>= 1 1)`, data.True)
	testCode(t, `(>= 0 1)`, data.False)
	testCode(t, `(>= 1 0)`, data.True)

	testCode(t, `(< 1 1)`, data.False)
	testCode(t, `(< 2 1)`, data.False)
	testCode(t, `(< 1 2)`, data.True)
	testCode(t, `(< 1 2 3 4 5)`, data.True)
	testCode(t, `(< 5 4 3 2 1)`, data.False)
	testCode(t, `(<= 1 1)`, data.True)
	testCode(t, `(<= 0 1)`, data.True)
	testCode(t, `(<= 1 0)`, data.False)
}

func TestBadCompare(t *testing.T) {
	e := interfaceErr("data.String", "data.Number", "Add")
	testBadCode(t, `(< 99 "hello")`, e)
	testBadCode(t, `(< "hello" "there")`, e)
}
