package builtin_test

import (
	"errors"
	"testing"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
)

func TestCondEval(t *testing.T) {
	as := assert.New(t)
	as.EvalTo(`(cond)`, data.Nil)

	as.EvalTo(`
		(cond
			[false "goodbye"]
			['()   "nope"]
			[true  "hello"]
			["hi"  "ignored"])
	`, S("hello"))

	as.EvalTo(`
		(cond
			[false "goodbye"]
			['()   "nope"]
			[:else "hello"]
			["hi"  "ignored"])
	`, S("hello"))

	as.EvalTo(`
		(cond
			[false "goodbye"]
			['()   "nope"])
	`, data.Nil)
}

func TestBadCond(t *testing.T) {
	as := assert.New(t)

	as.PanicWith(`
		(cond
			[true "hello"]
			[99])
	`, errors.New("invalid cond clause: [99]"))

	as.PanicWith(`
		(cond 99)
	`, errors.New("invalid cond clause: 99"))

	as.PanicWith(`
		(cond
			false "hello"
			99)
	`, errors.New("invalid cond clause: false"))
}
