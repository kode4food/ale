package assert

import (
	"github.com/kode4food/ale/pkg/data"
	"github.com/kode4food/ale/pkg/eval"
)

// Eval will evaluate source code on behalf of the test framework
func (w *Wrapper) Eval(src string) data.Value {
	w.Helper()
	ns := GetTestNamespace()
	res, err := eval.String(ns, data.String(src))
	if err != nil {
		panic(err)
	}
	return res
}

// EvalTo will evaluate source code and test for an expected result
func (w *Wrapper) EvalTo(src string, expect data.Value) {
	w.Helper()
	w.Equal(expect, w.Eval(src))
}

// PanicWith evaluates source code and expects a panic to happen
func (w *Wrapper) PanicWith(src string, err any) {
	w.Helper()
	defer w.ExpectPanic(err)
	w.Eval(src)
}
