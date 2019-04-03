package test

import (
	"fmt"
	"testing"

	"gitlab.com/kode4food/ale/api"
	. "gitlab.com/kode4food/ale/internal/assert/helpers"
)

func TestBasicNumber(t *testing.T) {
	testCode(t, `(+)`, F(0))
	testCode(t, `(*)`, F(1))
	testCode(t, `(+ 1 1)`, F(2.0))
	testCode(t, `(* 4 4)`, F(16.0))
	testCode(t, `(+ 5 4)`, F(9.0))
	testCode(t, `(* 12 3)`, F(36.0))
	testCode(t, `(- 10 4)`, F(6.0))
	testCode(t, `(- 10 4 2)`, F(4.0))
	testCode(t, `(/ 10 2)`, F(5.0))
	testCode(t, `(/ 10 2 5)`, F(1.0))
	testCode(t, `(mod 10 3)`, F(1.0))
	testCode(t, `(mod 100 8 7)`, F(4.0))
}

func TestNestedNumber(t *testing.T) {
	testCode(t, `(/ 10 (- 5 3))`, F(5.0))
	testCode(t, `(* 5 (- 5 3))`, F(10.0))
	testCode(t, `(/ 10 (/ 6 3))`, F(5.0))
}

func TestNonNumber(t *testing.T) {
	e := intfErr("api.String", "api.Number", "Add")
	testBadCode(t, `(+ 99 "hello")`, e)
	testBadCode(t, `(+ "hello")`, e)
}

func TestBadNumbers(t *testing.T) {
	testBadNumber := func(err string, ns string) {
		testBadCode(t, ns, fmt.Errorf(err, api.String(ns)))
	}

	testBadNumber(api.ExpectedInteger, "0xfkk")
	testBadNumber(api.ExpectedInteger, "0b01109")
	testBadNumber(api.ExpectedInteger, "123j-k")
	testBadNumber(api.ExpectedFloat, "1.2j-k")
	//testBadNumber(api.ExpectedRatio, "1/2p")
}

func TestCompare(t *testing.T) {
	testCode(t, `(= 1 1)`, api.True)
	testCode(t, `(= 1 1 1 1 '1 1 1)`, api.True)
	testCode(t, `(= 1 2)`, api.False)
	testCode(t, `(= 1 1 1 1 2 1 1 1)`, api.False)

	testCode(t, `(!= 1 1)`, api.False)
	testCode(t, `(!= 1 1 1 1 '1 1 1)`, api.False)
	testCode(t, `(!= 1 2)`, api.True)
	testCode(t, `(!= 1 1 1 1 2 1 1 1)`, api.True)

	testCode(t, `(> 1 1)`, api.False)
	testCode(t, `(> 2 1)`, api.True)
	testCode(t, `(> 1 2)`, api.False)
	testCode(t, `(> 1 2 3 4 5)`, api.False)
	testCode(t, `(> 5 4 3 2 1)`, api.True)
	testCode(t, `(>= 1 1)`, api.True)
	testCode(t, `(>= 0 1)`, api.False)
	testCode(t, `(>= 1 0)`, api.True)

	testCode(t, `(< 1 1)`, api.False)
	testCode(t, `(< 2 1)`, api.False)
	testCode(t, `(< 1 2)`, api.True)
	testCode(t, `(< 1 2 3 4 5)`, api.True)
	testCode(t, `(< 5 4 3 2 1)`, api.False)
	testCode(t, `(<= 1 1)`, api.True)
	testCode(t, `(<= 0 1)`, api.True)
	testCode(t, `(<= 1 0)`, api.False)
}

func TestBadCompare(t *testing.T) {
	e := intfErr("api.String", "api.Number", "Add")
	testBadCode(t, `(< 99 "hello")`, e)
	testBadCode(t, `(< "hello" "there")`, e)
}
