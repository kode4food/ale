package builtin_test

import (
	"fmt"
	"testing"

	"github.com/kode4food/ale/core/builtin"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/runtime/isa"
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

	str := "not a procedure"
	defer as.ExpectPanic(fmt.Errorf(builtin.ErrProcedureRequired, str))
	_ = builtin.Macro.Call(S(str))
}

func TestMacroExpand(t *testing.T) {
	testMacroExpandWith(t, builtin.MacroExpand)
}

func TestMacroExpand1(t *testing.T) {
	testMacroExpandWith(t, builtin.MacroExpand1)
}

func testMacroExpandWith(t *testing.T, enc testEncoder) {
	as := assert.New(t)
	e1 := assert.GetTestEncoder()

	neq := L(LS("declare"), LS("some-sym"))
	enc(e1, neq)
	e1.Emit(isa.Return)

	c := e1.Constants()
	as.Equal(2, len(c))
	s, ok := c[0].(data.Local)
	as.True(ok)
	as.Equal("some-sym", s.String())
	f, ok := c[1].(data.Procedure)
	as.True(ok)
	as.Equal("(ale/declare* some-sym)", data.ToString(f.Call(neq)))
}
