package generate

import (
	"gitlab.com/kode4food/ale/compiler/encoder"
	"gitlab.com/kode4food/ale/data"
	"gitlab.com/kode4food/ale/namespace"
	"gitlab.com/kode4food/ale/runtime/isa"
)

// Symbol encodes a symbol retrieval
func Symbol(e encoder.Type, s data.Symbol) {
	if l, ok := s.(data.LocalSymbol); ok {
		resolveLocal(e, l)
		return
	}
	resolveGlobal(e, s)
}

func resolveLocal(e encoder.Type, l data.LocalSymbol) {
	if scope, ok := e.ResolveScope(l); ok {
		switch scope {
		case encoder.LocalScope:
			idx, _ := e.ResolveLocal(l)
			e.Emit(isa.Load, idx)
		case encoder.ArgScope:
			idx, rest, _ := e.ResolveArg(l)
			if rest {
				e.Emit(isa.RestArg, idx)
			} else {
				e.Emit(isa.Arg, idx)
			}
		case encoder.NameScope:
			e.Emit(isa.Self)
		case encoder.ClosureScope:
			idx, _ := e.ResolveClosure(l)
			e.Emit(isa.Closure, idx)
		default:
			panic("unknown scope type")
		}
		return
	}
	resolveGlobal(e, l)
}

func resolveGlobal(e encoder.Type, s data.Symbol) {
	globals := e.Globals()
	entry := namespace.MustResolveSymbol(globals, s)
	if entry.IsBound() && entry.Owner() == globals {
		Literal(e, entry.Value())
		return
	}
	Literal(e, s)
	e.Emit(isa.Resolve)
}
