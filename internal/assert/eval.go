package assert

import (
	"gitlab.com/kode4food/ale/bootstrap"
	"gitlab.com/kode4food/ale/data"
	"gitlab.com/kode4food/ale/eval"
	"gitlab.com/kode4food/ale/namespace"
)

var (
	manager = namespace.NewManager()
	ready   bool
)

func (w *Wrapper) Eval(src string) data.Value {
	if !ready {
		bootstrap.Into(manager)
		ready = true
	}
	ns := manager.GetAnonymous()
	return eval.String(ns, data.String(src))
}

func (w *Wrapper) EvalTo(src string, expect data.Value) {
	w.Equal(expect, w.Eval(src))
}

func (w *Wrapper) PanicWith(src string, err error) {
	defer w.ExpectPanic(err.Error())
	w.Eval(src)
}
