package builtin_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/internal/runtime"
	"github.com/kode4food/ale/pkg/data"
)

func unexpectedTypeError(got, expected string) error {
	return fmt.Errorf(runtime.ErrUnexpectedType, got, expected)
}

func TestListEval(t *testing.T) {
	as := assert.New(t)
	as.MustEvalTo(`(list? '(1 2 3))`, data.True)
	as.MustEvalTo(`(list? '())`, data.True)
	as.MustEvalTo(`(list? [1 2 3])`, data.False)
	as.MustEvalTo(`(list? 42)`, data.False)
	as.MustEvalTo(`(list? (list 1 2 3))`, data.True)
	as.MustEvalTo(`(list)`, data.Null)

	as.MustEvalTo(`
		(define x '(1 2 3 4))
		(x 2)
	`, L(I(3), I(4)))
}

func TestVectorEval(t *testing.T) {
	as := assert.New(t)

	r1 := as.MustEval(`(vector 1 (- 5 3) (+ 1 2))`)
	as.String("[1 2 3]", r1)

	r2 := as.MustEval(`(apply vector (concat '(1) '((- 5 3)) '((+ 1 2))))`)
	as.String("[1 (- 5 3) (+ 1 2)]", r2)

	as.MustEvalTo(`(conj [1 2 3] 4)`, S("[1 2 3 4]"))
	as.MustEvalTo(`(vector? (conj [1 2 3] 4))`, data.True)

	as.MustEvalTo(`(vector? [1 2 3])`, data.True)
	as.MustEvalTo(`(vector? (vector 1 2 3))`, data.True)
	as.MustEvalTo(`(vector? [])`, data.True)
	as.MustEvalTo(`(vector? 99)`, data.False)

	as.MustEvalTo(`(!vector? [1 2 3])`, data.False)
	as.MustEvalTo(`(!vector? (vector 1 2 3))`, data.False)
	as.MustEvalTo(`(!vector? [])`, data.False)
	as.MustEvalTo(`(!vector? 99)`, data.True)

	as.MustEvalTo(`(counted? [1 2 3])`, data.True)
	as.MustEvalTo(`(counted? 99)`, data.False)
	as.MustEvalTo(`(indexed? [1 2 3])`, data.True)
	as.MustEvalTo(`(indexed? 99)`, data.False)

	as.MustEvalTo(`
		(define x [1 2 3 4])
		(x 2)
	`, V(I(3), I(4)))
}

func TestSequencesEval(t *testing.T) {
	as := assert.New(t)
	as.MustEvalTo(`(seq? [1 2 3])`, data.True)
	as.MustEvalTo(`(seq? '())`, data.True)
	as.MustEvalTo(`(empty? '())`, data.True)
	as.MustEvalTo(`(empty? '(1))`, data.False)
	as.MustEvalTo(`(seq '())`, data.False)
	as.MustEvalTo(`(seq? 99)`, data.False)
	as.MustEvalTo(`(seq 99)`, data.False)
}

func TestToObjectEval(t *testing.T) {
	as := assert.New(t)
	as.MustEvalTo(`(object? (seq->object [:name "Ale" :age 45]))`, data.True)
	as.MustEvalTo(`(object? (seq->object '(:name "Ale" :age 45)))`, data.True)
	as.MustEvalTo(`(mapped? (seq->object '(:name "Ale" :age 45)))`, data.True)
}

func TestToVectorEval(t *testing.T) {
	as := assert.New(t)
	as.MustEvalTo(`(vector? (seq->vector (list 1 2 3)))`, data.True)
	as.MustEvalTo(`(counted? [1 2 3 4])`, data.True)
	as.MustEvalTo(`(nth [1 2 3 4] 2)`, I(3))
	as.MustEvalTo(`(nth [1 2 3 4] 10 "oops")`, S("oops"))

	as.PanicWith(`(nth [1 2 3 4] 10)`, errors.New("index out of bounds"))
}

func TestToListEval(t *testing.T) {
	as := assert.New(t)
	as.MustEvalTo(`(list? (seq->list (vector 1 2 3)))`, data.True)
}

func TestMapFilterEval(t *testing.T) {
	as := assert.New(t)
	as.MustEvalTo(`
		(car (apply list (map (lambda (x) (* x 2)) [1 2 3 4])))
	`, F(2))

	as.MustEvalTo(`
		(define x (concat '(1 2) (list 3 4)))
		(define y
			(map
				(lambda (x) (* x 2))
				(filter
					(lambda (x) (= x 6))
					[5 6])))
		(apply +
			(map
				(lambda (z) (car z))
				[x y]))
	`, F(13))
}

func TestLenEval(t *testing.T) {
	as := assert.New(t)
	as.MustEvalTo(`
	  (length [1 2 3 4 5])
	`, I(5))

	as.MustEvalTo(`
		(length! (take 10000 (range 1 1000000000)))
	`, I(10000))

	as.PanicWith(`
		(length (take 10000 (range 1 1000000000)))
	`, unexpectedTypeError("lazy sequence", "counted"))
}

func TestLastEval(t *testing.T) {
	as := assert.New(t)
	as.MustEvalTo(`
	  (last [1 2 3 4 37])
	`, I(37))

	as.MustEvalTo(`
		(last! (take 10000 (range 1 1000000000 2)))
	`, I(19999))

	err := unexpectedTypeError("lazy sequence", "counted")
	as.PanicWith(`
		(last (take 10000 (range 1 1000000000)))
	`, err)
}

func TestReverse(t *testing.T) {
	as := assert.New(t)

	as.String(`(4 3 2 1)`, as.MustEval(`(reverse '(1 2 3 4))`))
	as.String(`[4 3 2 1]`, as.MustEval(`(reverse [1 2 3 4])`))
	as.MustEvalTo(`(reverse '())`, data.Null)
	as.MustEvalTo(`(reverse [])`, data.EmptyVector)
	as.String(`(4 3 2 1)`, as.MustEval(`(reverse! (take 4 (range 1 1000)))`))

	err := unexpectedTypeError("lazy sequence", "reverser")
	as.PanicWith(`
		(reverse (take 4 (range 1 1000)))
	`, err)
}
