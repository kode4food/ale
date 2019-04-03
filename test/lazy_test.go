package test

import (
	"testing"

	. "gitlab.com/kode4food/ale/internal/assert/helpers"
)

func TestRange(t *testing.T) {
	testCode(t, `
		(reduce
			(fn [x y] (+ x y))
			(range 1 5 1))
	`, F(10))

	testCode(t, `
		(reduce
			(fn [x y] (+ x y))
			(range 5 1 -1))
	`, F(14))
}

func TestMapAndFilter(t *testing.T) {
	testCode(t, `
		(reduce
			(fn [x y] (+ x y))
			(map
				(fn [x] (* x 2))
				(filter
					(fn [x] (<= x 5))
					[1 2 3 4 5 6 7 8 9 10])))
	`, F(30))
}

func TestMapParallel(t *testing.T) {
	testCode(t, `
		(to-vector
			(map +
				[1 2 3 4]
				'(2 4 6 8)
				(range 20 30)))
	`, S("[23 27 31 35]"))
}

func TestReduce(t *testing.T) {
	testCode(t, `
		(def x '(1 2 3 4))
		(reduce + x)
	`, F(10))

	testCode(t, `
		(def y (concat '(1 2 3 4) [5 6 7 8]))
		(reduce + y)
	`, F(36))

	testCode(t, `
		(def y (concat '(1 2 3 4) [5 6 7 8]))
		(reduce + 10 y)
	`, F(46))
}

func TestTakeDrop(t *testing.T) {
	testCode(t, `
		(def x (concat '(1 2 3 4) [5 6 7 8]))
		(nth (apply vector (take 6 x)) 5)
	`, F(6))

	testCode(t, `
		(def x (concat '(1 2 3 4) [5 6 7 8]))
		(nth (apply vector (drop 3 x)) 0)
	`, F(4))
}

func TestLazySeq(t *testing.T) {
	testCode(t, `
		(reduce
			(fn [x y] (+ x y))
			(lazy-seq (cons 1 (lazy-seq [2, 3]))))
	`, F(6))

	testCode(t, `
		(len (to-vector (lazy-seq nil)))
	`, F(0))
}

func TestForEachLoop(t *testing.T) {
	testCode(t, `
		(let [ch (chan) emit (:emit ch) close (:close ch) seq (:seq ch)]
			(go
				(for-each [i (range 1 5 1), j (range 1 10 2)]
					(emit (* i j)))
				(close))
			(reduce (fn [x y] (+ x y)) seq))
	`, F(250))
}
