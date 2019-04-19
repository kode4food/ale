package test

import (
	"testing"

	"gitlab.com/kode4food/ale/data"
	. "gitlab.com/kode4food/ale/internal/assert/helpers"
)

func TestSequences(t *testing.T) {
	testCode(t, `(seq? [1 2 3])`, data.True)
	testCode(t, `(seq? ())`, data.True)
	testCode(t, `(seq ())`, data.Nil)
	testCode(t, `(seq? 99)`, data.False)
	testCode(t, `(seq 99)`, data.Nil)
}

func TestToAssoc(t *testing.T) {
	testCode(t, `(assoc? (to-assoc [:name "Ale" :age 45]))`, data.True)
	testCode(t, `(assoc? (to-assoc '(:name "Ale" :age 45)))`, data.True)
	testCode(t, `(mapped? (to-assoc '(:name "Ale" :age 45)))`, data.True)
}

func TestToVector(t *testing.T) {
	testCode(t, `(vector? (to-vector (list 1 2 3)))`, data.True)
	testCode(t, `(len? [1 2 3 4])`, data.True)
}

func TestToList(t *testing.T) {
	testCode(t, `(list? (to-list (vector 1 2 3)))`, data.True)
}

func TestMapFilter(t *testing.T) {
	testCode(t, `
		(first (apply list (map (fn [x] (* x 2)) [1 2 3 4])))
	`, F(2))

	testCode(t, `
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
