package eval_test

import (
	"testing"

	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/pkg/core/bootstrap"
	"github.com/kode4food/ale/pkg/data"
	"github.com/kode4food/ale/pkg/eval"
	"github.com/kode4food/ale/pkg/read"
)

func TestBasicEval(t *testing.T) {
	as := assert.New(t)

	e := bootstrap.DevNullEnvironment()
	ns := e.GetAnonymous()

	v1 := eval.String(ns, "(if true 1 0)")
	as.Number(1, v1)

	v2 := eval.String(ns, "((lambda (x) (* x 2)) 50)")
	as.Number(100, v2)

	v3 := eval.String(ns, "(first (concat [1 2 3] [4 5 6]))")
	as.Number(1, v3)

	eval.String(ns, "(define x 99)")
	x, ok := ns.Resolve("x")
	as.True(ok && x.IsBound())
	as.Number(99, x.Value())

	v4 := eval.String(ns, "(and true true)")
	as.True(v4)
}

func TestBuiltIns(t *testing.T) {
	as := assert.New(t)

	e := bootstrap.DevNullEnvironment()
	b := e.GetAnonymous()
	ns := e.GetRoot()

	ns.Declare("hello").Bind(
		data.MakeProcedure(func(...data.Value) data.Value {
			return S("there")
		}, 0),
	)

	tr := read.FromString(`(hello)`)
	as.String("there", eval.Block(b, tr))
}
