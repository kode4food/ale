package macro

import (
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/env"
	"github.com/kode4food/ale/internal/sequence"
)

// Expand performs a complete macro expansion
func Expand(ns env.Namespace, v data.Value) data.Value {
	if res, ok := expand1(ns, v); ok {
		return Expand(ns, res)
	}
	return v
}

// Expand1 performs a single macro expansion
func Expand1(ns env.Namespace, v data.Value) data.Value {
	res, _ := expand1(ns, v)
	return res
}

func expand1(ns env.Namespace, v data.Value) (data.Value, bool) {
	l, ok := v.(*data.List) // it's got to be a list
	if !ok {
		return v, false
	}
	f, r, _ := l.Split() // starting with a symbol
	s, ok := f.(data.Symbol)
	if !ok {
		return v, false
	}
	args := sequence.ToValues(r)
	rv, ok := env.ResolveValue(ns, s) // that actually resolves
	if !ok {
		return v, false
	}
	m, ok := rv.(Call) // to a macro call
	if !ok {
		return v, false
	}
	return m(ns, args...), true
}
