package assert

import (
	"github.com/kode4food/ale/pkg/data"
	"github.com/kode4food/ale/pkg/env"
)

func (w *Wrapper) IsNotDeclared(ns env.Namespace, n data.Local) {
	w.Helper()
	e, err := ns.Resolve(n)
	w.Nil(e)
	w.NotNil(err)
}

func (w *Wrapper) IsNotBound(ns env.Namespace, n data.Local) {
	w.Helper()
	e, err := ns.Resolve(n)
	w.NotNil(e)
	w.Nil(err)
	w.False(e.IsBound())
	v, err := e.Value()
	w.Nil(v)
	w.NotNil(err)
}

func (w *Wrapper) IsBound(ns env.Namespace, n data.Local) data.Value {
	w.Helper()
	e, err := ns.Resolve(n)
	w.NotNil(e)
	w.Nil(err)
	w.True(e.IsBound())
	v, err := e.Value()
	w.Nil(err)
	return v
}
