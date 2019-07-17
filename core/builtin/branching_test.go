package builtin_test

import (
	"errors"
	"testing"

	"gitlab.com/kode4food/ale/data"
	"gitlab.com/kode4food/ale/internal/assert"
	. "gitlab.com/kode4food/ale/internal/assert/helpers"
)

func TestCondEval(t *testing.T) {
	as := assert.New(t)
	as.EvalTo(`(cond)`, data.Null)

	as.EvalTo(`
		(cond
			[#f   "goodbye"]
			['()  "nope"]
			[#t   "hello"]
			["hi" "ignored"])
	`, S("hello"))

	as.EvalTo(`
		(cond
			[#f    "goodbye"]
			['()   "nope"]
			[:else "hello"]
			["hi"  "ignored"])
	`, S("hello"))

	as.EvalTo(`
		(cond
			[#f  "goodbye"]
			['() "nope"])
	`, data.Null)
}

func TestBadCond(t *testing.T) {
	as := assert.New(t)

	as.PanicWith(`
		(cond
			[#t "hello"]
			[99])
	`, errors.New("cond clause must be paired: [99]"))

	as.PanicWith(`
		(cond 99)
	`, errors.New("cond clause must be a vector: 99"))

	as.PanicWith(`
		(cond
			#f "hello"
			99)
	`, errors.New("cond clause must be a vector: #f"))
}
