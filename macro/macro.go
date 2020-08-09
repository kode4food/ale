package macro

import (
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/sequence"
	"github.com/kode4food/ale/namespace"
)

// Expand performs a complete macro expansion
func Expand(ns namespace.Type, v data.Value) data.Value {
	if res, ok := expand1(ns, v); ok {
		return Expand(ns, res)
	}
	return v
}

// Expand1 performs a single macro expansion
func Expand1(ns namespace.Type, v data.Value) data.Value {
	res, _ := expand1(ns, v)
	return res
}

func expand1(ns namespace.Type, v data.Value) (data.Value, bool) {
	if l, ok := v.(data.List); ok {
		if s, ok := l.First().(data.Symbol); ok {
			args := sequence.ToValues(l.Rest())
			if v, ok := namespace.ResolveValue(ns, s); ok {
				if m, ok := v.(Call); ok {
					return m(ns, args...), true
				}
			}
		}
	}
	return v, false
}
