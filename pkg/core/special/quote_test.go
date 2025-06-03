package special_test

import (
	"testing"

	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/pkg/core/bootstrap"
	"github.com/kode4food/ale/pkg/data"
	"github.com/kode4food/ale/pkg/env"
	"github.com/kode4food/ale/pkg/eval"
)

func TestQuoteEval(t *testing.T) {
	as := assert.New(t)

	r1 := as.MustEval("(quote (blah 2 3))").(*data.List)
	r2 := as.MustEval("'(blah 2 3)").(*data.List)

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

	e := bootstrap.DevNullEnvironment()
	ns := e.GetAnonymous()

	as.Nil(env.BindPublic(ns, "foo", F(456)))
	r1, err := eval.String(ns, `'[123 ,foo]`)
	as.NoError(err)
	as.String("[123 (ale/unquote foo)]", r1)
}

func TestUnquoteMacroEval(t *testing.T) {
	as := assert.New(t)
	as.MustEvalTo(
		"(define-macro test (x . y) `(,x ,@y {:hello 99}))"+
			"(test vector 1 2 3)",
		S("[1 2 3 {:hello 99}]"),
	)
}
