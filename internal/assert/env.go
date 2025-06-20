package assert

import (
	"github.com/kode4food/ale/pkg/data"
	"github.com/kode4food/ale/pkg/env"
)

func (w *Wrapper) IsNotDeclared(ns env.Namespace, n data.Local) {
	w.Helper()
	e, in, err := ns.Resolve(n)
	w.Nil(e)
	w.Nil(in)
	w.NotNil(err)
}

func (w *Wrapper) IsNotBound(ns env.Namespace, n data.Local) {
	w.Helper()
	e, in, err := ns.Resolve(n)
	if w.NoError(err) {
		w.NotNil(e)
		w.NotNil(in)
		w.False(e.IsBound())

		v, err := e.Value()
		w.Nil(v)
		w.NotNil(err)
	}
}

func (w *Wrapper) IsBound(ns env.Namespace, n data.Local) data.Value {
	w.Helper()
	e, in, err := ns.Resolve(n)
	if w.NoError(err) {
		w.NotNil(e)
		w.NotNil(in)

		w.True(e.IsBound())
		v, err := e.Value()
		w.NoError(err)
		return v
	}
	return nil
}
