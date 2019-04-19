package test

import (
	"testing"

	"gitlab.com/kode4food/ale/data"
	. "gitlab.com/kode4food/ale/internal/assert/helpers"
)

func TestCond(t *testing.T) {
	testCode(t, `(cond)`, data.Nil)

	testCode(t, `
		(cond
			false "goodbye"
			nil   "nope"
			true  "hello"
			"hi"  "ignored")
	`, S("hello"))

	testCode(t, `
		(cond
			false "goodbye"
			nil   "nope"
			:else "hello"
			"hi"  "ignored")
	`, S("hello"))

	testCode(t, `
		(cond
			false "goodbye"
			nil   "nope")
	`, data.Nil)

	testCode(t, `
		(cond
			true "hello"
			99)
	`, S("hello"))

	testCode(t, `(cond 99)`, F(99))

	testCode(t, `
		(cond
			false "hello"
			99)
	`, F(99))
}
