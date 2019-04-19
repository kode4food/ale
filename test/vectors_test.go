package test

import (
	"testing"

	"gitlab.com/kode4food/ale/data"
	"gitlab.com/kode4food/ale/internal/assert"
	. "gitlab.com/kode4food/ale/internal/assert/helpers"
)

func TestVector(t *testing.T) {
	as := assert.New(t)

	r1 := runCode(`(vector 1 (- 5 3) (+ 1 2))`)
	as.String("[1 2 3]", r1)

	r2 := runCode(`(apply vector (concat '(1) '((- 5 3)) '((+ 1 2))))`)
	as.String("[1 (- 5 3) (+ 1 2)]", r2)

	testCode(t, `(conj [1 2 3] 4)`, S("[1 2 3 4]"))
	testCode(t, `(vector? (conj [1 2 3] 4))`, data.True)

	testCode(t, `(vector? [1 2 3])`, data.True)
	testCode(t, `(vector? (vector 1 2 3))`, data.True)
	testCode(t, `(vector? [])`, data.True)
	testCode(t, `(vector? 99)`, data.False)

	testCode(t, `(!vector? [1 2 3])`, data.False)
	testCode(t, `(!vector? (vector 1 2 3))`, data.False)
	testCode(t, `(!vector? [])`, data.False)
	testCode(t, `(!vector? 99)`, data.True)

	testCode(t, `(len? [1 2 3])`, data.True)
	testCode(t, `(len? 99)`, data.False)
	testCode(t, `(indexed? [1 2 3])`, data.True)
	testCode(t, `(indexed? 99)`, data.False)

	testCode(t, `
		(def x [1 2 3 4])
		(x 2)
	`, F(3))
}
