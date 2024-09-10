package core_test

import (
	"fmt"
	"testing"

	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/pkg/core"
	"github.com/kode4food/ale/pkg/core/internal"
	"github.com/kode4food/ale/pkg/data"
)

func getPredicate(kwd data.Keyword) data.Procedure {
	return core.IsA.Call(kwd).(data.Procedure)
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

	as.PanicWith(`(null?)`, fmt.Errorf(internal.ErrUnmatchedCase, 0, "1 or more"))

	as.PanicWith(`(is-a :dog "woof!")`,
		fmt.Errorf(core.ErrUnknownPredicate, data.Keyword("dog")),
	)
}

func TestTypeOf(t *testing.T) {
	as := assert.New(t)
	as.EvalTo(`
		(define l-pred (type-of* '(1 2 3)))
		[(l-pred '(9 8 7))
         (l-pred '())
         (l-pred [1 2 3])
         (eq l-pred (type-of* '(9 8 7)))
		 (eq l-pred (type-of* []))
         (eq l-pred (type-of* '()))]
	`, V(data.True, data.False, data.False, data.True, data.False, data.False))
}
