package builtin_test

import (
	"fmt"
	"testing"

	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/pkg/core/builtin"
	"github.com/kode4food/ale/pkg/data"
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
			(foo1 (args 0) (args 1) (args 2)))

		(macroexpand '(foo2 1 2 3))
	`, S("hello123"))
}

func TestBadMacro(t *testing.T) {
	as := assert.New(t)

	defer as.ExpectPanic(fmt.Errorf(builtin.ErrProcedureRequired, "32"))
	_ = builtin.Macro.Call(I(32))
}
