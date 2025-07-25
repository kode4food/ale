package builtin_test

import (
	"fmt"
	"testing"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/internal/runtime"
)

func TestBasicNumberEval(t *testing.T) {
	as := assert.New(t)
	as.MustEvalTo(`(+)`, I(0))
	as.MustEvalTo(`(*)`, I(1))
	as.MustEvalTo(`(+ 1 1)`, I(2))
	as.MustEvalTo(`(* 4 4)`, I(16))
	as.MustEvalTo(`(+ 5 4)`, I(9))
	as.MustEvalTo(`(* 12 3)`, I(36))
	as.MustEvalTo(`(- 10 4)`, I(6))
	as.MustEvalTo(`(- 10 4 2)`, I(4))
	as.MustEvalTo(`(- 5)`, I(-5))
	as.MustEvalTo(`(/ 10 2)`, I(5))
	as.MustEvalTo(`(/ 10 2 5)`, I(1))
	as.MustEvalTo(`(mod 10 3)`, I(1))
	as.MustEvalTo(`(mod 100 8 7)`, I(4))
}

func TestNestedNumberEval(t *testing.T) {
	as := assert.New(t)
	as.MustEvalTo(`(/ 10 (- 5 3))`, I(5))
	as.MustEvalTo(`(* 5 (- 5 3))`, I(10))
	as.MustEvalTo(`(/ 10 (/ 6 3))`, I(5))
}

func TestNonNumberEval(t *testing.T) {
	as := assert.New(t)
	as.MustEvalTo(`(inf? (/ 99.0 0))`, data.True)
	as.MustEvalTo(`(inf? 99)`, data.False)

	as.MustEvalTo(`(-inf? (/ -99.0 0))`, data.True)
	as.MustEvalTo(`(-inf? -99)`, data.False)

	as.MustEvalTo(`(nan? 99)`, data.False)
	as.MustEvalTo(`(nan? (/ 0.0 0))`, data.True)
	as.MustEvalTo(`(nan? "hello")`, data.False)
}

func TestBadMathsEval(t *testing.T) {
	as := assert.New(t)
	err := fmt.Errorf(runtime.ErrUnexpectedType, "string", "number")

	as.PanicWith(`(+ 99 "hello")`, err)
	as.PanicWith(`(+ "hello")`, err)
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
	as.MustEvalTo(`(= 1 1)`, data.True)
	as.MustEvalTo(`(= 1 1 1 1 '1 1 1)`, data.True)
	as.MustEvalTo(`(= 1 2)`, data.False)
	as.MustEvalTo(`(= 1 1 1 1 2 1 1 1)`, data.False)

	as.MustEvalTo(`(!= 1 1)`, data.False)
	as.MustEvalTo(`(!= 1 1 1 1 '1 1 1)`, data.False)
	as.MustEvalTo(`(!= 1 2)`, data.True)
	as.MustEvalTo(`(!= 1 1 1 1 2 1 1 1)`, data.True)

	as.MustEvalTo(`(> 1 1)`, data.False)
	as.MustEvalTo(`(> 2 1)`, data.True)
	as.MustEvalTo(`(> 1 2)`, data.False)
	as.MustEvalTo(`(> 1 2 3 4 5)`, data.False)
	as.MustEvalTo(`(> 5 4 3 2 1)`, data.True)
	as.MustEvalTo(`(>= 1 1)`, data.True)
	as.MustEvalTo(`(>= 0 1)`, data.False)
	as.MustEvalTo(`(>= 1 0)`, data.True)

	as.MustEvalTo(`(< 1 1)`, data.False)
	as.MustEvalTo(`(< 2 1)`, data.False)
	as.MustEvalTo(`(< 1 2)`, data.True)
	as.MustEvalTo(`(< 1 2 3 4 5)`, data.True)
	as.MustEvalTo(`(< 5 4 3 2 1)`, data.False)
	as.MustEvalTo(`(<= 1 1)`, data.True)
	as.MustEvalTo(`(<= 0 1)`, data.True)
	as.MustEvalTo(`(<= 1 0)`, data.False)
}

func TestBadCompareEval(t *testing.T) {
	as := assert.New(t)
	err := fmt.Errorf(runtime.ErrUnexpectedType, "string", "number")
	as.PanicWith(`(< 99 "hello")`, err)
	as.PanicWith(`(< "hello" "there")`, err)
}
