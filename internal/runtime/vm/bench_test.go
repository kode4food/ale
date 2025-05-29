package vm_test

import (
	"testing"

	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/internal/runtime/isa"
	"github.com/kode4food/ale/pkg/core/bootstrap"
	"github.com/kode4food/ale/pkg/data"
	"github.com/kode4food/ale/pkg/eval"
)

func numExplicitSum(n1, n2, n3 data.Integer) data.Value {
	return n1.Add(n2).Add(n3)
}

func numLoopSum(args ...data.Value) data.Value {
	var sum data.Number = I(0)
	for _, a := range args {
		sum = sum.Add(a.(data.Integer))
	}
	return sum
}

func nativeExplicitSum(i1, i2, i3 int) int {
	return i1 + i2 + i3
}

func nativeLoopSum(args ...int) int {
	var sum = 0
	for _, v := range args {
		sum += v
	}
	return sum
}

func BenchmarkNativeExplicit(b *testing.B) {
	for range b.N {
		_ = nativeExplicitSum(5, 6, 7) + 12
	}
}

func BenchmarkNativeLoop(b *testing.B) {
	for range b.N {
		_ = nativeLoopSum(5, 6, 7) + 12
	}
}

func BenchmarkNumberExplicit(b *testing.B) {
	i1 := I(5)
	i2 := I(6)
	i3 := I(7)
	i4 := I(12)
	for range b.N {
		_ = numExplicitSum(i1, i2, i3).(data.Integer) + i4
	}
}

func BenchmarkNumberLoop(b *testing.B) {
	i1 := I(5)
	i2 := I(6)
	i3 := I(7)
	i4 := I(12)
	for range b.N {
		_ = numLoopSum(i1, i2, i3).(data.Integer) + i4
	}
}

var bCode = makeClosure(isa.Instructions{
	isa.Const.New(0), // the extra stack item is intentional
	isa.Const.New(0),
	isa.Const.New(1),
	isa.PosInt.New(1),
	isa.Const.New(3),
	isa.Call.New(3),
	isa.Add.New(),
	isa.Return.New(),
})

func BenchmarkCalls(b *testing.B) {
	for range b.N {
		_ = bCode.Call()
	}
}

func BenchmarkBottles(b *testing.B) {
	env := bootstrap.DevNullEnvironment()
	ns := env.GetAnonymous()
	res, err := eval.String(ns, `
		(define println no-op) ; bypass the OS
		(define (bottles n)
		  (str
		    (cond [(= n 0) "No more bottles"]
		          [(= n 1) "One bottle"]
		          [:else (str n " bottles")])
		    " of beer"))
		(define (beer n)
		  (when (> n 0)
		    (println (bottles n) "on the wall")
		    (println (bottles n)) (println "Take one down, pass it around")
		    (println (bottles (- n 1)) "on the wall")
		    (println)
		    (beer (- n 1))))
		beer
	`)
	if err != nil {
		panic(err)
	}
	beer := res.(data.Procedure)

	b.ResetTimer()
	for range b.N {
		_ = beer.Call(data.Integer(99))
	}
}

func BenchmarkTailCalls(b *testing.B) {
	env := bootstrap.DevNullEnvironment()
	ns := env.GetAnonymous()
	res, _ := eval.String(ns, `
		(define (fib-iter curr prev count)
			(if (= count 0)
				prev
				(fib-iter (+ curr prev) curr (dec count))))
		(define (fib n)
			(fib-iter 1 0 n))
		fib
	`)
	fib := res.(data.Procedure)

	b.ResetTimer()
	for range b.N {
		res := fib.Call(data.Integer(10))
		if res, ok := res.(data.Integer); !ok || !res.Equal(data.Integer(55)) {
			b.Fatalf("fib result is incorrect")
		}
	}
}

func BenchmarkNonTailCalls(b *testing.B) {
	env := bootstrap.DevNullEnvironment()
	ns := env.GetAnonymous()
	res, _ := eval.String(ns, `
		(define (fib n)
			(cond
				[(= n 0) 0]
				[(= n 1) 1]
				[:else (+ (fib (- n 1)) (fib (- n 2)))]))
		fib
	`)
	fib := res.(data.Procedure)

	b.ResetTimer()
	for range b.N {
		res := fib.Call(data.Integer(10))
		if res, ok := res.(data.Integer); !ok || !res.Equal(data.Integer(55)) {
			b.Fatalf("fib result is incorrect")
		}
	}
}
