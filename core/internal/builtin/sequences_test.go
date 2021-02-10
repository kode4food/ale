package builtin_test

import (
	"errors"
	"testing"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
)

func TestSequencesEval(t *testing.T) {
	as := assert.New(t)
	as.EvalTo(`(seq? [1 2 3])`, data.True)
	as.EvalTo(`(seq? '())`, data.True)
	as.EvalTo(`(empty? '())`, data.True)
	as.EvalTo(`(empty? '(1))`, data.False)
	as.EvalTo(`(seq '())`, data.Nil)
	as.EvalTo(`(seq? 99)`, data.False)
	as.EvalTo(`(seq 99)`, data.Nil)
}

func TestToObjectEval(t *testing.T) {
	as := assert.New(t)
	as.EvalTo(`(object? (seq->object [:name "Ale" :age 45]))`, data.True)
	as.EvalTo(`(object? (seq->object '(:name "Ale" :age 45)))`, data.True)
	as.EvalTo(`(mapped? (seq->object '(:name "Ale" :age 45)))`, data.True)
}

func TestToVectorEval(t *testing.T) {
	as := assert.New(t)
	as.EvalTo(`(vector? (seq->vector (list 1 2 3)))`, data.True)
	as.EvalTo(`(counted? [1 2 3 4])`, data.True)
	as.EvalTo(`(nth [1 2 3 4] 2)`, data.Integer(3))
	as.EvalTo(`(nth [1 2 3 4] 10 "oops")`, data.String("oops"))

	as.PanicWith(`(nth [1 2 3 4] 10)`, errors.New("index out of bounds"))
}

func TestToListEval(t *testing.T) {
	as := assert.New(t)
	as.EvalTo(`(list? (seq->list (vector 1 2 3)))`, data.True)
}

func TestMapFilterEval(t *testing.T) {
	as := assert.New(t)
	as.EvalTo(`
		(first (apply list (map (lambda (x) (* x 2)) [1 2 3 4])))
	`, F(2))

	as.EvalTo(`
		(define x (concat '(1 2) (list 3 4)))
		(define y
			(map
				(lambda (x) (* x 2))
				(filter
					(lambda (x) (= x 6))
					[5 6])))
		(apply +
			(map
				(lambda (z) (first z))
				[x y]))
	`, F(13))
}

func TestLenEval(t *testing.T) {
	as := assert.New(t)
	as.EvalTo(`
	  (length [1 2 3 4 5])
	`, I(5))

	as.EvalTo(`
		(length! (take 10000 (range 1 1000000000)))
	`, I(10000))

	as.PanicWith(`
		(length (take 10000 (range 1 1000000000)))
	`, interfaceErr("*sequence.lazySequence", "data.CountedSequence"))
}

func TestLastEval(t *testing.T) {
	as := assert.New(t)
	as.EvalTo(`
	  (last [1 2 3 4 37])
	`, I(37))

	as.EvalTo(`
		(last! (take 10000 (range 1 1000000000 2)))
	`, I(19999))

	err := interfaceErr("*sequence.lazySequence", "data.CountedSequence")
	as.PanicWith(`
		(last (take 10000 (range 1 1000000000)))
	`, err)
}

func TestReverse(t *testing.T) {
	as := assert.New(t)

	as.String(`(4 3 2 1)`, as.Eval(`(reverse '(1 2 3 4))`))
	as.String(`[4 3 2 1]`, as.Eval(`(reverse [1 2 3 4])`))
	as.EvalTo(`(reverse '())`, data.EmptyList)
	as.EvalTo(`(reverse [])`, data.EmptyVector)
	as.String(`(4 3 2 1)`, as.Eval(`(reverse! (take 4 (range 1 1000)))`))

	err := interfaceErr("*sequence.lazySequence", "data.ReverserSequence")
	as.PanicWith(`
		(reverse (take 4 (range 1 1000)))
	`, err)
}
