package builtin_test

import (
	"fmt"
	"testing"

	"gitlab.com/kode4food/ale/compiler/arity"
	"gitlab.com/kode4food/ale/data"
	"gitlab.com/kode4food/ale/internal/assert"
)

func TestPredicatesEval(t *testing.T) {
	as := assert.New(t)
	as.EvalTo(`(eq #t #t #t)`, data.True)
	as.EvalTo(`(eq #t #f #t)`, data.False)
	as.EvalTo(`(eq #f #f #f)`, data.True)

	as.EvalTo(`(!eq #t #t #t)`, data.False)
	as.EvalTo(`(!eq #t #f)`, data.True)
	as.EvalTo(`(!eq #f #f)`, data.False)

	as.EvalTo(`(null? '())`, data.True)
	as.EvalTo(`(null? '() '() '())`, data.True)
	as.EvalTo(`(null? #f)`, data.False)
	as.EvalTo(`(null? #f '())`, data.False)

	as.EvalTo(`(null? "hello")`, data.False)
	as.EvalTo(`(null? '(1 2 3))`, data.False)
	as.EvalTo(`(null? '() "hello")`, data.False)

	as.EvalTo(`(keyword? :hello)`, data.True)
	as.EvalTo(`(!keyword? :hello)`, data.False)
	as.EvalTo(`(keyword? 99)`, data.False)
	as.EvalTo(`(!keyword? 99)`, data.True)

	as.PanicWith(`(null?)`, fmt.Errorf(arity.BadMinimumArity, 0, 1))
}
