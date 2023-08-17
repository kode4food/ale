package builtin_test

import (
	"fmt"
	"testing"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
)

func TestBasicNumberEval(t *testing.T) {
	as := assert.New(t)
	as.EvalTo(`(+)`, I(0))
	as.EvalTo(`(*)`, I(1))
	as.EvalTo(`(+ 1 1)`, I(2))
	as.EvalTo(`(* 4 4)`, I(16))
	as.EvalTo(`(+ 5 4)`, I(9))
	as.EvalTo(`(* 12 3)`, I(36))
	as.EvalTo(`(- 10 4)`, I(6))
	as.EvalTo(`(- 10 4 2)`, I(4))
	as.EvalTo(`(- 5)`, I(-5))
	as.EvalTo(`(/ 10 2)`, I(5))
	as.EvalTo(`(/ 10 2 5)`, I(1))
	as.EvalTo(`(mod 10 3)`, I(1))
	as.EvalTo(`(mod 100 8 7)`, I(4))
}

func TestNestedNumberEval(t *testing.T) {
	as := assert.New(t)
	as.EvalTo(`(/ 10 (- 5 3))`, I(5))
	as.EvalTo(`(* 5 (- 5 3))`, I(10))
	as.EvalTo(`(/ 10 (/ 6 3))`, I(5))
}

func TestNonNumberEval(t *testing.T) {
	as := assert.New(t)
	as.EvalTo(`(is-pos-inf (/ 99.0 0))`, data.True)
	as.EvalTo(`(is-pos-inf 99)`, data.False)
	as.EvalTo(`(is-pos-inf "hello")`, data.False)

	as.EvalTo(`(is-neg-inf (/ -99.0 0))`, data.True)
	as.EvalTo(`(is-neg-inf -99)`, data.False)
	as.EvalTo(`(is-neg-inf "hello")`, data.False)

	as.EvalTo(`(is-nan 99)`, data.False)
	as.EvalTo(`(is-nan "hello")`, data.False)
}

func TestBadMathsEval(t *testing.T) {
	as := assert.New(t)
	e := unexpectedTypeError("string", "number")

	as.PanicWith(`(+ 99 "hello")`, e)
	as.PanicWith(`(+ "hello")`, e)
}

func TestBadNumbersEval(t *testing.T) {
	as := assert.New(t)

	testBadNumber := func(err string, ns string) {
		as.PanicWith(ns, fmt.Errorf(err, S(ns)))
	}

	testBadNumber(data.ErrExpectedInteger, "0xfkk")
	testBadNumber(data.ErrExpectedInteger, "0b01109")
	testBadNumber(data.ErrExpectedInteger, "123j-k")
	testBadNumber(data.ErrExpectedFloat, "1.2j-k")
	testBadNumber(data.ErrExpectedRatio, "1/2p")
}

func TestCompareEval(t *testing.T) {
	as := assert.New(t)
	as.EvalTo(`(= 1 1)`, data.True)
	as.EvalTo(`(= 1 1 1 1 '1 1 1)`, data.True)
	as.EvalTo(`(= 1 2)`, data.False)
	as.EvalTo(`(= 1 1 1 1 2 1 1 1)`, data.False)

	as.EvalTo(`(!= 1 1)`, data.False)
	as.EvalTo(`(!= 1 1 1 1 '1 1 1)`, data.False)
	as.EvalTo(`(!= 1 2)`, data.True)
	as.EvalTo(`(!= 1 1 1 1 2 1 1 1)`, data.True)

	as.EvalTo(`(> 1 1)`, data.False)
	as.EvalTo(`(> 2 1)`, data.True)
	as.EvalTo(`(> 1 2)`, data.False)
	as.EvalTo(`(> 1 2 3 4 5)`, data.False)
	as.EvalTo(`(> 5 4 3 2 1)`, data.True)
	as.EvalTo(`(>= 1 1)`, data.True)
	as.EvalTo(`(>= 0 1)`, data.False)
	as.EvalTo(`(>= 1 0)`, data.True)

	as.EvalTo(`(< 1 1)`, data.False)
	as.EvalTo(`(< 2 1)`, data.False)
	as.EvalTo(`(< 1 2)`, data.True)
	as.EvalTo(`(< 1 2 3 4 5)`, data.True)
	as.EvalTo(`(< 5 4 3 2 1)`, data.False)
	as.EvalTo(`(<= 1 1)`, data.True)
	as.EvalTo(`(<= 0 1)`, data.True)
	as.EvalTo(`(<= 1 0)`, data.False)
}

func TestBadCompareEval(t *testing.T) {
	as := assert.New(t)
	e := unexpectedTypeError("string", "number")
	as.PanicWith(`(< 99 "hello")`, e)
	as.PanicWith(`(< "hello" "there")`, e)
}
