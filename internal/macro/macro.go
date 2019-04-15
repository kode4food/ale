package macro

import (
	"gitlab.com/kode4food/ale/api"
	"gitlab.com/kode4food/ale/internal/namespace"
	"gitlab.com/kode4food/ale/stdlib"
)

// Expand performs a complete macro expansion
func Expand(ns namespace.Type, v api.Value) api.Value {
	if res, ok := expand1(ns, v); ok {
		return Expand(ns, res)
	}
	return v
}

// Expand1 performs a single macro expansion
func Expand1(ns namespace.Type, v api.Value) api.Value {
	res, _ := expand1(ns, v)
	return res
}

func expand1(ns namespace.Type, v api.Value) (api.Value, bool) {
	if l, ok := v.(*api.List); ok {
		if s, ok := l.First().(api.Symbol); ok {
			args := stdlib.SequenceToVector(l.Rest())
			if v, ok := namespace.ResolveSymbol(ns, s); ok {
				if m, ok := v.(*api.Function); ok && m.IsMacro() {
					return m.Call(args...), true
				}
			}
			if s == syntaxSym {
				return SyntaxQuote(ns, args[0]), true
			}
		}
	}
	return v, false
}
