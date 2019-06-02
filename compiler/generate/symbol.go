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
	n := l.Name()
	if s, ok := e.ResolveScoped(n); ok {
		switch s.Scope {
		case encoder.LocalScope:
			c, _ := e.ResolveLocal(n)
			e.Emit(isa.Load, c.Index)
			if c.Type == encoder.ReferenceCell {
				e.Emit(isa.Deref)
			}
		case encoder.ArgScope:
			c, _ := e.ResolveArg(n)
			if c.Type == encoder.RestCell {
				e.Emit(isa.RestArg, c.Index)
			} else {
				e.Emit(isa.Arg, c.Index)
			}
		case encoder.NameScope:
			e.Emit(isa.Self)
		case encoder.ClosureScope:
			c, _ := e.ResolveClosure(n)
			e.Emit(isa.Closure, c.Index)
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
