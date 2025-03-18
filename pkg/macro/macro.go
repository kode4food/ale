package macro

import (
	"github.com/kode4food/ale/internal/sequence"
	"github.com/kode4food/ale/pkg/data"
	"github.com/kode4food/ale/pkg/env"
)

// Expand performs a complete macro expansion
func Expand(ns env.Namespace, v data.Value) (data.Value, error) {
	if res, ok := expand1(ns, v); ok {
		return Expand(ns, res)
	}
	return v, nil
}

// Expand1 performs a single macro expansion
func Expand1(ns env.Namespace, v data.Value) (data.Value, error) {
	res, _ := expand1(ns, v)
	return res, nil
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
	rv, err := env.ResolveValue(ns, s) // that actually resolves
	if err != nil {
		return v, false
	}
	m, ok := rv.(Call) // to a macro call
	if !ok {
		return v, false
	}
	args := sequence.ToVector(r)
	return m(ns, args...), true
}
