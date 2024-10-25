package builtin_test

import (
	"testing"

	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/pkg/core/builtin"
	"github.com/kode4food/ale/pkg/data"
)

func TestLazySequence(t *testing.T) {
	as := assert.New(t)

	var i int
	var p data.Procedure

	p = data.MakeProcedure(func(...data.Value) data.Value {
		if i < 10 {
			p := builtin.LazySequence.Call(p).(data.Prepender)
			res := p.Prepend(data.Integer(i))
			i++
			return res
		}
		return data.Null
	}, 0)

	s := builtin.LazySequence.Call(p).(data.Sequence)
	as.String(`(0 1 2 3 4 5 6 7 8 9)`, data.MakeSequenceStr(s))
}

func TestRangeEval(t *testing.T) {
	as := assert.New(t)
	as.MustEvalTo(`
		(fold-left
			(lambda (x y) (+ x y))
			(range 1 5 1))
	`, F(10))

	as.MustEvalTo(`
		(fold-left
			(lambda (x y) (+ x y))
			(range 5 1 -1))
	`, F(14))
}

func TestMapAndFilterEval(t *testing.T) {
	as := assert.New(t)
	as.MustEvalTo(`
		(fold-left
			(lambda (x y) (+ x y))
			(map
				(lambda (x) (* x 2))
				(filter
					(lambda (x) (<= x 5))
					[1 2 3 4 5 6 7 8 9 10])))
	`, F(30))
}

func TestMapParallelEval(t *testing.T) {
	as := assert.New(t)
	as.MustEvalTo(`
		(seq->vector
			(map +
				[1 2 3 4]
				'(2 4 6 8)
				(range 20 30)))
	`, S("[23 27 31 35]"))
}

func TestReduceEval(t *testing.T) {
	as := assert.New(t)
	as.MustEvalTo(`
		(define x '(1 2 3 4))
		(fold-left + x)
	`, F(10))

	as.MustEvalTo(`
		(define y (concat '(1 2 3 4) [5 6 7 8]))
		(fold-left + y)
	`, F(36))

	as.MustEvalTo(`
		(define y (concat '(1 2 3 4) [5 6 7 8]))
		(fold-left + 10 y)
	`, F(46))
}

func TestTakeDropEval(t *testing.T) {
	as := assert.New(t)
	as.MustEvalTo(`
		(define x (concat '(1 2 3 4) [5 6 7 8]))
		(nth (apply vector (take 6 x)) 5)
	`, F(6))

	as.MustEvalTo(`
		(define x (concat '(1 2 3 4) [5 6 7 8]))
		(nth (apply vector (drop 3 x)) 0)
	`, F(4))

	as.PanicWith(
		`(last! (drop 99 57))`,
		unexpectedTypeError("integer", "pair"),
	)
	as.PanicWith(
		`(last! (take 99 57))`,
		unexpectedTypeError("integer", "sequence"),
	)
}

func TestLazySeqEval(t *testing.T) {
	as := assert.New(t)
	as.MustEvalTo(`
		(fold-left
			(lambda (x y) (+ x y))
			(lazy-seq (cons 1 (lazy-seq [2 3]))))
	`, F(6))

	as.MustEvalTo(`
		(length (seq->vector (lazy-seq '())))
	`, F(0))
}

func TestForEachLoopEval(t *testing.T) {
	as := assert.New(t)
	as.MustEvalTo(`
		(let* ([ch (chan)]
			   [emit (:emit ch)]
			   [close (:close ch)]
			   [seq (:seq ch)])
			(go
				(for-each ([i (range 1 5 1)]
				           [j (range 1 10 2)])
					(emit (* i j)))
				(close))
			(fold-left (lambda (x y) (+ x y)) seq))
	`, F(250))
}
