package builtin_test

import (
	"fmt"
	"testing"

	"github.com/kode4food/ale/compiler/arity"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/assert"
)

func TestPredicatesEval(t *testing.T) {
	as := assert.New(t)
	as.EvalTo(`(eq true true true)`, data.True)
	as.EvalTo(`(eq true false true)`, data.False)
	as.EvalTo(`(eq false false false)`, data.True)

	as.EvalTo(`(!eq true true true)`, data.False)
	as.EvalTo(`(!eq true false)`, data.True)
	as.EvalTo(`(!eq false false)`, data.False)

	as.EvalTo(`(null? '())`, data.True)
	as.EvalTo(`(null? '() '() '())`, data.True)
	as.EvalTo(`(null? false)`, data.False)
	as.EvalTo(`(null? false '())`, data.False)

	as.EvalTo(`(null? "hello")`, data.False)
	as.EvalTo(`(null? '(1 2 3))`, data.False)
	as.EvalTo(`(null? '() "hello")`, data.False)

	as.EvalTo(`(keyword? :hello)`, data.True)
	as.EvalTo(`(!keyword? :hello)`, data.False)
	as.EvalTo(`(keyword? 99)`, data.False)
	as.EvalTo(`(!keyword? 99)`, data.True)

	as.PanicWith(`(null?)`, fmt.Errorf(arity.BadMinimumArity, 0, 1))
}
