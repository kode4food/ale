package eval_test

import (
	"testing"

	"gitlab.com/kode4food/ale/bootstrap"
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

	v1 := eval.String(ns, "(if true 1 0)")
	as.Integer(1, v1)

	v2 := eval.String(ns, "((fn [x] (* x 2)) 50)")
	as.Integer(100, v2)

	v3 := eval.String(ns, "(first (concat [1 2 3] [4 5 6]))")
	as.Integer(1, v3)

	eval.String(ns, "(def x 99)")
	v4, ok := ns.Resolve(data.Name("x"))
	as.True(ok)
	as.Integer(99, v4)

	v5 := eval.String(ns, "(and true true)")
	as.True(v5)
}

func TestBuiltIns(t *testing.T) {
	as := assert.New(t)

	manager := namespace.NewManager()
	bootstrap.Into(manager)
	b := manager.GetAnonymous()
	ns := manager.GetRoot()

	ns.Bind("hello", data.ApplicativeFunction(func(_ ...data.Value) data.Value {
		return S("there")
	}))

	l := read.Scan(`(hello)`)
	tr := read.FromScanner(l)

	as.String("there", eval.Block(b, tr))
}
