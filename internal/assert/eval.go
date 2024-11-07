package assert

import (
	"github.com/kode4food/ale/pkg/data"
	"github.com/kode4food/ale/pkg/eval"
)

// MustEval will evaluate source code on behalf of the test framework
func (w *Wrapper) MustEval(src string) data.Value {
	w.Helper()
	res, err := w.Eval(src)
	if err != nil {
		panic(err)
	}
	return res
}

func (w *Wrapper) Eval(src string) (data.Value, error) {
	w.Helper()
	ns := GetTestNamespace()
	return eval.String(ns, data.String(src))
}

// MustEvalTo will evaluate source code and test for an expected result
func (w *Wrapper) MustEvalTo(src string, expect data.Value) {
	w.Helper()
	w.Equal(expect, w.MustEval(src))
}

func (w *Wrapper) ErrorWith(src string, err any) {
	w.Helper()
	_, res := w.Eval(src)
	w.ExpectError(err, res)
}

// PanicWith evaluates source code and expects a panic to happen
func (w *Wrapper) PanicWith(src string, err any) {
	w.Helper()
	defer w.ExpectPanic(err)
	w.MustEval(src)
}
