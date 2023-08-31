package vm_test

import (
	"testing"

	"github.com/kode4food/ale/core/bootstrap"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/eval"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/runtime/isa"
)

func numExplicitSum(n1, n2, n3 data.Integer) data.Value {
	return n1 + n2 + n3
}

func numLoopSum(args ...data.Value) data.Value {
	var sum = I(0)
	for _, a := range args {
		sum = sum + a.(data.Integer)
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
	for n := 0; n < b.N; n++ {
		_ = nativeExplicitSum(5, 6, 7) + 12
	}
}

func BenchmarkNativeLoop(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_ = nativeLoopSum(5, 6, 7) + 12
	}
}

func BenchmarkNumberExplicit(b *testing.B) {
	i1 := I(5)
	i2 := I(6)
	i3 := I(7)
	i4 := I(12)
	for n := 0; n < b.N; n++ {
		_ = numExplicitSum(i1, i2, i3).(data.Integer) + i4
	}
}

func BenchmarkNumberLoop(b *testing.B) {
	i1 := I(5)
	i2 := I(6)
	i3 := I(7)
	i4 := I(12)
	for n := 0; n < b.N; n++ {
		_ = numLoopSum(i1, i2, i3).(data.Integer) + i4
	}
}

var bCode = makeFunc(isa.Instructions{
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
	for n := 0; n < b.N; n++ {
		_ = bCode.Call()
	}
}

func BenchmarkTailCalls(b *testing.B) {
	env := bootstrap.DevNullEnvironment()
	ns := env.GetAnonymous()
	_ = eval.String(ns, `
		(define (bottles n)
		  (str
		    (cond [ (= n 0) "No more bottles"]
		          [ (= n 1) "One bottle"]
		          [:else (str n " bottles")])
		    " of beer"))
		(define (beer n)
		  (when (> n 0)
		    (println (bottles n) "on the wall")
		    (println (bottles n)) (println "Take one down, pass it around")
		    (println (bottles (- n 1)) "on the wall")
		    (println)
		    (beer (- n 1))))
	`)
	entry, _ := ns.Resolve("beer")
	beer := entry.Value().(data.Function)
	for n := 0; n < b.N; n++ {
		_ = beer.Call(data.Integer(9999))
	}
}
