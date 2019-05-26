package builtin_test

import (
	"testing"

	"gitlab.com/kode4food/ale/data"
	"gitlab.com/kode4food/ale/internal/assert"
	. "gitlab.com/kode4food/ale/internal/assert/helpers"
)

func TestMacroPredicatesEval(t *testing.T) {
	as := assert.New(t)
	as.EvalTo(`(macro? cond)`, data.True)
	as.EvalTo(`(!macro? cond)`, data.False)
	as.EvalTo(`(macro? if)`, data.False)
	as.EvalTo(`(!macro? if)`, data.True)
}

func TestMacroReplaceEval(t *testing.T) {
	as := assert.New(t)
	as.EvalTo(`
		(defmacro foo
			[& args]
			(to-list (cons 'str (cons "hello" args))))

		(foo 1 2 3)
	`, S(`hello123`))
}

func TestMacroExpandEval(t *testing.T) {
	as := assert.New(t)
	as.EvalTo(`
		(defmacro foo1
			[& args]
			(to-list (cons 'str (cons "hello" args))))

		(macroexpand-1 '(foo1 1 2 3))
	`, S(`(str "hello" 1 2 3)`))

	as.EvalTo(`
		(defmacro foo1
			[& args]
			(to-list (cons 'str (cons "hello" args))))

		(defmacro foo2
			[& args]
			(foo1 (args 0) (args 1) (args 2)))

		(macroexpand '(foo2 1 2 3))
	`, S("hello123"))
}
