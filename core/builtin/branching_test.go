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
			false "goodbye"
			null  "nope"
			true  "hello"
			"hi"  "ignored")
	`, S("hello"))

	as.EvalTo(`
		(cond
			false "goodbye"
			null  "nope"
			:else "hello"
			"hi"  "ignored")
	`, S("hello"))

	as.EvalTo(`
		(cond
			false "goodbye"
			null  "nope")
	`, data.Null)

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
