package builtin_test

import (
	"fmt"
	"testing"

	"github.com/kode4food/ale/core/internal/builtin"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/assert"
)

func getPredicate(kwd data.Keyword) data.Function {
	return builtin.IsA.Call(kwd).(data.Function)
}

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

	as.EvalTo(
		`(eq { :name "Ale" :age 3 :colors [:red :green :blue] }
             { :age 3 :colors [:red :green :blue] :name "Ale" })`,
		data.True,
	)

	as.PanicWith(`(null?)`, fmt.Errorf(data.ErrMinimumArity, 1, 0))

	as.PanicWith(`(is-a* :dog "woof!")`,
		fmt.Errorf(builtin.ErrUnknownPredicate, data.Keyword("dog")),
	)
}
