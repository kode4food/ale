package vm_test

import (
	"testing"

	"github.com/kode4food/ale/data"
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

var bCode = makeCode([]isa.Coder{
	isa.Const, isa.Index(0), // the extra stack item is intentional
	isa.Const, isa.Index(0),
	isa.Const, isa.Index(1),
	isa.One,
	isa.Const, isa.Index(3),
	isa.Call, isa.Count(3),
	isa.Add,
	isa.Return,
})

func BenchmarkVMCalls(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_ = bCode.Call()
	}
}
