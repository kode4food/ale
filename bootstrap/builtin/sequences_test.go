package builtin_test

import (
	"testing"

	"gitlab.com/kode4food/ale/data"
	"gitlab.com/kode4food/ale/internal/assert"
	. "gitlab.com/kode4food/ale/internal/assert/helpers"
)

func TestSequencesEval(t *testing.T) {
	as := assert.New(t)
	as.EvalTo(`(seq? [1 2 3])`, data.True)
	as.EvalTo(`(seq? ())`, data.True)
	as.EvalTo(`(empty? ())`, data.True)
	as.EvalTo(`(empty? '(1))`, data.False)
	as.EvalTo(`(seq ())`, data.Nil)
	as.EvalTo(`(seq? 99)`, data.False)
	as.EvalTo(`(seq 99)`, data.Nil)
}

func TestToAssocEval(t *testing.T) {
	as := assert.New(t)
	as.EvalTo(`(assoc? (to-assoc [:name "Ale" :age 45]))`, data.True)
	as.EvalTo(`(assoc? (to-assoc '(:name "Ale" :age 45)))`, data.True)
	as.EvalTo(`(mapped? (to-assoc '(:name "Ale" :age 45)))`, data.True)
}

func TestToVectorEval(t *testing.T) {
	as := assert.New(t)
	as.EvalTo(`(vector? (to-vector (list 1 2 3)))`, data.True)
	as.EvalTo(`(len? [1 2 3 4])`, data.True)
}

func TestToListEval(t *testing.T) {
	as := assert.New(t)
	as.EvalTo(`(list? (to-list (vector 1 2 3)))`, data.True)
}

func TestMapFilterEval(t *testing.T) {
	as := assert.New(t)
	as.EvalTo(`
		(first (apply list (map (fn [x] (* x 2)) [1 2 3 4])))
	`, F(2))

	as.EvalTo(`
		(def x (concat '(1 2) (list 3 4)))
		(def y
			(map
				(fn [x] (* x 2))
				(filter
					(fn [x] (= x 6))
					[5 6])))
		(apply +
			(map
				(fn [z] (first z))
				[x y]))
	`, F(13))
}
