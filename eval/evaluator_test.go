package eval_test

import (
	"testing"

	"gitlab.com/kode4food/ale/core/bootstrap"
	"gitlab.com/kode4food/ale/data"
	"gitlab.com/kode4food/ale/eval"
	"gitlab.com/kode4food/ale/internal/assert"
	. "gitlab.com/kode4food/ale/internal/assert/helpers"
	"gitlab.com/kode4food/ale/namespace"
	"gitlab.com/kode4food/ale/read"
)

func TestBasicEval(t *testing.T) {
	as := assert.New(t)

	manager := namespace.NewManager()
	bootstrap.Into(manager)
	ns := manager.GetAnonymous()

	v1 := eval.String(ns, "(if #t 1 0)")
	as.Number(1, v1)

	v2 := eval.String(ns, "((fn [x] (* x 2)) 50)")
	as.Number(100, v2)

	v3 := eval.String(ns, "(first (concat [1 2 3] [4 5 6]))")
	as.Number(1, v3)

	eval.String(ns, "(def x 99)")
	e, ok := ns.Resolve(data.Name("x"))
	as.True(ok && e.IsBound())
	as.Number(99, e.Value())

	v4 := eval.String(ns, "(and #t #t)")
	as.True(v4)
}

func TestBuiltIns(t *testing.T) {
	as := assert.New(t)

	manager := namespace.NewManager()
	bootstrap.Into(manager)
	b := manager.GetAnonymous()
	ns := manager.GetRoot()

	ns.Declare("hello").Bind(
		data.Call(func(_ ...data.Value) data.Value {
			return S("there")
		}),
	)

	l := read.Scan(`(hello)`)
	tr := read.FromScanner(l)

	as.String("there", eval.Block(b, tr))
}
