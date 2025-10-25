package builtin_test

import (
	"fmt"
	"testing"

	"github.com/kode4food/ale"
	"github.com/kode4food/ale/core/builtin"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/macro"
)

func TestMacroPredicatesEval(t *testing.T) {
	as := assert.New(t)
	as.MustEvalTo(`(macro? cond)`, data.True)
	as.MustEvalTo(`(!macro? cond)`, data.False)
	as.MustEvalTo(`(macro? if)`, data.False)
	as.MustEvalTo(`(!macro? if)`, data.True)
	as.MustEvalTo(`(atom? "hello")`, data.True)
	as.MustEvalTo(`(!atom? '(1 2 3))`, data.True)
}

func TestMacroReplaceEval(t *testing.T) {
	as := assert.New(t)
	as.MustEvalTo(`
		(define-macro foo args
			(seq->list (cons 'str (cons "hello" args))))

		(foo 1 2 3)
	`, S(`hello123`))
}

func TestMacroExpandEval(t *testing.T) {
	as := assert.New(t)
	as.MustEvalTo(`
		(define-macro foo1 args
			(seq->list (cons 'str (cons "hello" args))))

		(macroexpand-1 '(foo1 1 2 3))
	`, S(`(str "hello" 1 2 3)`))

	as.MustEvalTo(`
		(define-macro foo1 args
			(seq->list (cons 'str (cons "hello" args))))

		(define-macro foo2 args
			(foo1 (0 args) (1 args) (2 args)))

		(macroexpand '(foo2 1 2 3))
	`, S("hello123"))
}

func TestBadMacro(t *testing.T) {
	as := assert.New(t)
	as.Panics(
		func() { _ = builtin.Macro.Call(F(32)) },
		fmt.Errorf("%w: %s", builtin.ErrProcedureRequired, "32"),
	)

	ns := assert.GetTestEnvironment().GetRoot()
	m := builtin.Macro.Call(data.MakeProcedure(
		func(...ale.Value) ale.Value {
			return data.False
		}, 1),
	).(macro.Call)

	as.Panics(
		func() { m(ns, I(1), I(2)) },
		fmt.Errorf(data.ErrFixedArity, 1, 2))
}
