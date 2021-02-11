package builtin_test

import (
	"testing"

	"github.com/kode4food/ale/core/bootstrap"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/env"
	"github.com/kode4food/ale/eval"
	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
)

func TestQuoteEval(t *testing.T) {
	as := assert.New(t)

	r1 := as.Eval("(quote (blah 2 3))").(data.List)
	r2 := as.Eval("'(blah 2 3)").(data.List)

	v1, ok := r1.ElementAt(0)
	v2, _ := r2.ElementAt(0)
	as.True(ok)
	as.Equal(v1, v2)

	v1, ok = r1.ElementAt(0)
	as.True(ok)
	if _, ok = v1.(data.Symbol); !ok {
		as.Fail("first element is not a symbol")
	}

	v1, ok = r1.ElementAt(1)
	v2, _ = r2.ElementAt(1)
	as.True(ok)
	as.Identical(v1, v2)

	v1, ok = r1.ElementAt(1)
	as.True(ok)
	as.Number(2, v1)

	v1, ok = r1.ElementAt(2)
	v2, _ = r2.ElementAt(2)
	as.True(ok)
	as.Identical(v1, v2)

	v1, ok = r1.ElementAt(2)
	as.True(ok)
	as.Number(3, v1)
}

func TestUnquoteEval(t *testing.T) {
	as := assert.New(t)

	e := env.NewEnvironment()
	bootstrap.Into(e)
	ns := e.GetAnonymous()

	ns.Declare("foo").Bind(F(456))
	r1 := eval.String(ns, `'[123 ,foo]`)
	as.String("[123 (ale/unquote foo)]", r1)
}

func TestUnquoteMacroEval(t *testing.T) {
	as := assert.New(t)
	as.EvalTo(
		"(define-macro test (x . y) `(,x ,@y {:hello 99}))"+
			"(test vector 1 2 3)",
		S("[1 2 3 {:hello 99}]"),
	)
}
