package builtin_test

import (
	"fmt"
	"testing"

	"github.com/kode4food/ale/core/internal/builtin"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
)

func TestMacroPredicatesEval(t *testing.T) {
	as := assert.New(t)
	as.EvalTo(`(macro? cond)`, data.True)
	as.EvalTo(`(!macro? cond)`, data.False)
	as.EvalTo(`(macro? if)`, data.False)
	as.EvalTo(`(!macro? if)`, data.True)
	as.EvalTo(`(atom? "hello")`, data.True)
	as.EvalTo(`(!atom? '(1 2 3))`, data.True)
}

func TestMacroReplaceEval(t *testing.T) {
	as := assert.New(t)
	as.EvalTo(`
		(define-macro foo args
			(seq->list (cons 'str (cons "hello" args))))

		(foo 1 2 3)
	`, S(`hello123`))
}

func TestMacroExpandEval(t *testing.T) {
	as := assert.New(t)
	as.EvalTo(`
		(define-macro foo1 args
			(seq->list (cons 'str (cons "hello" args))))

		(macroexpand-1 '(foo1 1 2 3))
	`, S(`(str "hello" 1 2 3)`))

	as.EvalTo(`
		(define-macro foo1 args
			(seq->list (cons 'str (cons "hello" args))))

		(define-macro foo2 args
			(foo1 (args 0) (args 1) (args 2)))

		(macroexpand '(foo2 1 2 3))
	`, S("hello123"))
}

func TestBadMacro(t *testing.T) {
	as := assert.New(t)

	str := "not a function"
	defer as.ExpectPanic(fmt.Errorf(builtin.ErrFunctionRequired, str))
	_ = builtin.Macro.Call(S(str))
}
