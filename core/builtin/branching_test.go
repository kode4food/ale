package builtin_test

import (
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
			#f   "goodbye"
			'()  "nope"
			#t   "hello"
			"hi" "ignored")
	`, S("hello"))

	as.EvalTo(`
		(cond
			#f    "goodbye"
			'()   "nope"
			:else "hello"
			"hi"  "ignored")
	`, S("hello"))

	as.EvalTo(`
		(cond
			#f  "goodbye"
			'() "nope")
	`, data.Null)

	as.EvalTo(`
		(cond
			#t "hello"
			99)
	`, S("hello"))

	as.EvalTo(`(cond 99)`, F(99))

	as.EvalTo(`
		(cond
			#f "hello"
			99)
	`, F(99))
}
