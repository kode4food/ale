package special_test

import (
	"errors"
	"testing"

	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/data"
)

func TestCondEval(t *testing.T) {
	as := assert.New(t)
	as.MustEvalTo(`(cond)`, data.Null)

	as.MustEvalTo(`
		(cond
			[false "goodbye"]
			[true  "hello"]
			["hi"  "ignored"])
	`, S("hello"))

	as.MustEvalTo(`
		(cond
			[false "goodbye"]
			[:else "hello"]
			["hi"  "ignored"])
	`, S("hello"))

	as.MustEvalTo(`
		(cond
			[false "goodbye"]
			['()   "hello"])
	`, S("hello"))
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
