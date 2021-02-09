package assert

import (
	"github.com/kode4food/ale/core/bootstrap"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/env"
	"github.com/kode4food/ale/eval"
)

var (
	testEnv = env.NewEnvironment()
	ready   bool
)

// Eval will evaluate source code on behalf of the test framework
func (w *Wrapper) Eval(src string) data.Value {
	w.Helper()
	if !ready {
		bootstrap.Into(testEnv)
		ready = true
	}
	ns := testEnv.GetAnonymous()
	return eval.String(ns, data.String(src))
}

// EvalTo will evaluate source code and test for an expected result
func (w *Wrapper) EvalTo(src string, expect data.Value) {
	w.Helper()
	w.Equal(expect, w.Eval(src))
}

// PanicWith will evaluate source code and expect a panic to happen
func (w *Wrapper) PanicWith(src string, err error) {
	w.Helper()
	defer w.ExpectPanic(err.Error())
	w.Eval(src)
}
