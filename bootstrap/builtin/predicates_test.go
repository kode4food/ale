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
	as.EvalTo(`(eq true true true)`, data.True)
	as.EvalTo(`(eq true false true)`, data.False)
	as.EvalTo(`(eq false false false)`, data.True)

	as.EvalTo(`(!eq true true true)`, data.False)
	as.EvalTo(`(!eq true false)`, data.True)
	as.EvalTo(`(!eq false false)`, data.False)

	as.EvalTo(`(nil? nil)`, data.True)
	as.EvalTo(`(nil? nil nil nil)`, data.True)
	as.EvalTo(`(nil? () nil)`, data.False)
	as.EvalTo(`(nil? false)`, data.False)
	as.EvalTo(`(nil? false () nil)`, data.False)

	as.EvalTo(`(nil? "hello")`, data.False)
	as.EvalTo(`(nil? '(1 2 3))`, data.False)
	as.EvalTo(`(nil? () nil "hello")`, data.False)

	as.EvalTo(`(keyword? :hello)`, data.True)
	as.EvalTo(`(!keyword? :hello)`, data.False)
	as.EvalTo(`(keyword? 99)`, data.False)
	as.EvalTo(`(!keyword? 99)`, data.True)

	as.PanicWith(`(nil?)`, fmt.Errorf(arity.BadMinimumArity, 0, 1))
}
