package builtin_test

import (
	"testing"

	"gitlab.com/kode4food/ale/data"
	"gitlab.com/kode4food/ale/internal/assert"
	. "gitlab.com/kode4food/ale/internal/assert/helpers"
)

func TestCondEval(t *testing.T) {
	as := assert.New(t)
	as.EvalTo(`(cond)`, data.Nil)

	as.EvalTo(`
		(cond
			false "goodbye"
			nil   "nope"
			true  "hello"
			"hi"  "ignored")
	`, S("hello"))

	as.EvalTo(`
		(cond
			false "goodbye"
			nil   "nope"
			:else "hello"
			"hi"  "ignored")
	`, S("hello"))

	as.EvalTo(`
		(cond
			false "goodbye"
			nil   "nope")
	`, data.Nil)

	as.EvalTo(`
		(cond
			true "hello"
			99)
	`, S("hello"))

	as.EvalTo(`(cond 99)`, F(99))

	as.EvalTo(`
		(cond
			false "hello"
			99)
	`, F(99))
}