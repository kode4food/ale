package builtin_test

import (
	"fmt"
	"testing"

	"github.com/kode4food/ale/core/builtin"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/internal/lang/params"
)

func getPredicate(kwd data.Keyword) data.Procedure {
	return builtin.IsA.Call(kwd).(data.Procedure)
}

func TestPredicatesEval(t *testing.T) {
	as := assert.New(t)
	as.MustEvalTo(`(eq true true true)`, data.True)
	as.MustEvalTo(`(eq true false true)`, data.False)
	as.MustEvalTo(`(eq false false false)`, data.True)

	as.MustEvalTo(`(!eq true true true)`, data.False)
	as.MustEvalTo(`(!eq true false)`, data.True)
	as.MustEvalTo(`(!eq false false)`, data.False)

	as.MustEvalTo(`(null? '())`, data.True)
	as.MustEvalTo(`(null? '() '() '())`, data.True)
	as.MustEvalTo(`(null? false)`, data.False)
	as.MustEvalTo(`(null? false '())`, data.False)

	as.MustEvalTo(`(null? "hello")`, data.False)
	as.MustEvalTo(`(null? '(1 2 3))`, data.False)
	as.MustEvalTo(`(null? '() "hello")`, data.False)

	as.MustEvalTo(`(keyword? :hello)`, data.True)
	as.MustEvalTo(`(!keyword? :hello)`, data.False)
	as.MustEvalTo(`(keyword? 99)`, data.False)
	as.MustEvalTo(`(!keyword? 99)`, data.True)

	as.MustEvalTo(
		`(eq { :name "Ale" :age 3 :colors [:red :green :blue] }
             { :age 3 :colors [:red :green :blue] :name "Ale" })`,
		data.True,
	)

	as.PanicWith(`(null?)`, fmt.Errorf(params.ErrUnmatchedCase, 0, "1 or more"))

	as.PanicWith(`(is-a :dog "woof!")`,
		fmt.Errorf("%w: %s", builtin.ErrUnknownPredicate, data.Keyword("dog")),
	)
}

func TestTypeOf(t *testing.T) {
	as := assert.New(t)
	as.MustEvalTo(`
		(define l-pred (%type-of '(1 2 3)))
		[(l-pred '(1 2 3))
		 (eq l-pred (%type-of '(1 2 3)))
		 (l-pred '(9 8 7))
         (l-pred '())
         (l-pred [1 2 3])
         (eq l-pred (%type-of '(9 8 7)))
		 (eq l-pred (%type-of []))
         (eq l-pred (%type-of '()))]
	`, V(
		data.True, data.True, data.False, data.False, data.False, data.False,
		data.False, data.False,
	))
}
