package vm_test

import (
	"testing"

	"github.com/kode4food/ale/internal/assert"
	"github.com/kode4food/ale/internal/runtime/vm"
	"github.com/kode4food/ale/pkg/data"
)

func TestProcedureHashCode(t *testing.T) {
	as := assert.New(t)
	p1 := as.Eval(`(lambda (x) (* x 2))`).(*vm.Closure)
	p2 := as.Eval(`(lambda (y) (* y 2))`).(*vm.Closure)
	p3 := as.Eval(`(lambda (x) (/ x 2))`).(*vm.Closure)
	as.True(p1.HashCode() == p2.HashCode())
	as.False(p1.HashCode() == p3.HashCode())
	as.False(p3.HashCode() == p2.HashCode())
}

func TestProcedureCaptured(t *testing.T) {
	as := assert.New(t)
	res := as.Eval(`
		(define (make op left) (lambda (x) (op left x)))
		[(make + 1) (make + 1) (make + 2) (make - 1)]
	`).(data.Vector)

	as.Equal(4, len(res))
	as.True(res[0].Equal(res[1]))
	as.True(data.HashCode(res[0]) == data.HashCode(res[1]))
	as.False(res[0].Equal(res[2]))
	as.False(data.HashCode(res[0]) == data.HashCode(res[2]))
	as.False(res[0].Equal(res[3]))
	as.False(data.HashCode(res[0]) == data.HashCode(res[3]))
}
