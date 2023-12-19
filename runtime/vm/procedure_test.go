package vm_test

import (
	"testing"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
)

func TestClosureEqual(t *testing.T) {
	as := assert.New(t)
	as.EvalTo(`
		(define f1 (lambda (x) (* x 2)))
		(define f2 (lambda (x) (* x 2)))
		(define f3 (lambda (x) (/ x 2)))
		(define f4 (lambda (y) (* y 2)))

		[(eq f1 f1) (eq f1 f2) (eq f1 f3) (eq f1 f4)]
	`, V(data.True, data.True, data.False, data.True))
}
