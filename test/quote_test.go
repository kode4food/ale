package test

import (
	"testing"

	"gitlab.com/kode4food/ale/api"
	"gitlab.com/kode4food/ale/eval"
	"gitlab.com/kode4food/ale/internal/assert"
	. "gitlab.com/kode4food/ale/internal/assert/helpers"
	"gitlab.com/kode4food/ale/internal/bootstrap"
)

func TestQuote(t *testing.T) {
	as := assert.New(t)

	r1 := runCode("(quote (blah 2 3))").(*api.List)
	r2 := runCode("'(blah 2 3)").(*api.List)

	v1, ok := r1.ElementAt(0)
	v2, _ := r2.ElementAt(0)
	as.True(ok)
	as.Equal(v1, v2)

	v1, ok = r1.ElementAt(0)
	as.True(ok)
	if _, ok = v1.(api.Symbol); !ok {
		as.Fail("first element is not a symbol")
	}

	v1, ok = r1.ElementAt(1)
	v2, _ = r2.ElementAt(1)
	as.True(ok)
	as.Identical(v1, v2)

	v1, ok = r1.ElementAt(1)
	as.True(ok)
	as.Integer(2, v1)

	v1, ok = r1.ElementAt(2)
	v2, _ = r2.ElementAt(2)
	as.True(ok)
	as.Identical(v1, v2)

	v1, ok = r1.ElementAt(2)
	as.True(ok)
	as.Integer(3, v1)
}

func TestUnquote(t *testing.T) {
	as := assert.New(t)

	manager := bootstrap.NullManager()
	bootstrap.Into(manager)
	ns := manager.GetUserNamespace()

	ns.Bind("foo", api.Float(456))
	r1 := eval.String(ns, `'[123 ~foo]`)
	as.String("[123 (ale/unquote foo)]", r1)
}

func TestUnquoteMacro(t *testing.T) {
	testCode(t,
		"(defmacro test [x & y] `(~x ~@y {:hello 99}))"+
			"(test vector 1 2 3)",
		S("[1 2 3 {:hello 99}]"),
	)
}
