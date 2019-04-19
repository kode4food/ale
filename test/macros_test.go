package test

import (
	"testing"

	"gitlab.com/kode4food/ale/data"
	. "gitlab.com/kode4food/ale/internal/assert/helpers"
)

func TestMacroPredicates(t *testing.T) {
	testCode(t, `(macro? cond)`, data.True)
	testCode(t, `(!macro? cond)`, data.False)
	testCode(t, `(macro? if)`, data.False)
	testCode(t, `(!macro? if)`, data.True)
}

func TestMacroReplace(t *testing.T) {
	testCode(t, `
		(defmacro foo
			[& args]
			(to-list (cons 'str (cons "hello" args))))

		(foo 1 2 3)
	`, S(`hello123`))
}

func TestMacroExpand(t *testing.T) {
	testCode(t, `
		(defmacro foo1
			[& args]
			(to-list (cons 'str (cons "hello" args))))

		(macroexpand-1 '(foo1 1 2 3))
	`, S(`(str "hello" 1 2 3)`))

	testCode(t, `
		(defmacro foo1
			[& args]
			(to-list (cons 'str (cons "hello" args))))

		(defmacro foo2
			[& args]
			(foo1 (args 0) (args 1) (args 2)))

		(macroexpand '(foo2 1 2 3))
	`, S("hello123"))
}
