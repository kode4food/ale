package builtin_test

import (
	"testing"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
)

func TestVectorEval(t *testing.T) {
	as := assert.New(t)

	r1 := as.Eval(`(vector 1 (- 5 3) (+ 1 2))`)
	as.String("[1 2 3]", r1)

	r2 := as.Eval(`(apply vector (concat '(1) '((- 5 3)) '((+ 1 2))))`)
	as.String("[1 (- 5 3) (+ 1 2)]", r2)

	as.EvalTo(`(conj [1 2 3] 4)`, S("[1 2 3 4]"))
	as.EvalTo(`(vector? (conj [1 2 3] 4))`, data.True)

	as.EvalTo(`(vector? [1 2 3])`, data.True)
	as.EvalTo(`(vector? (vector 1 2 3))`, data.True)
	as.EvalTo(`(vector? [])`, data.True)
	as.EvalTo(`(vector? 99)`, data.False)

	as.EvalTo(`(!vector? [1 2 3])`, data.False)
	as.EvalTo(`(!vector? (vector 1 2 3))`, data.False)
	as.EvalTo(`(!vector? [])`, data.False)
	as.EvalTo(`(!vector? 99)`, data.True)

	as.EvalTo(`(counted? [1 2 3])`, data.True)
	as.EvalTo(`(counted? 99)`, data.False)
	as.EvalTo(`(indexed? [1 2 3])`, data.True)
	as.EvalTo(`(indexed? 99)`, data.False)

	as.EvalTo(`
		(define x [1 2 3 4])
		(x 2)
	`, F(3))
}
