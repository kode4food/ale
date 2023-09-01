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
	if l, ok := v.(data.List); ok {
		f, r, _ := l.Split()
		if s, ok := f.(data.Symbol); ok {
			args := sequence.ToValues(r)
			if v, ok := env.ResolveValue(ns, s); ok {
				if m, ok := v.(Call); ok {
					return m(ns, args...), true
				}
			}
		}
	}
	return v, false
}
