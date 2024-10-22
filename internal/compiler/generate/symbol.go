package generate

import (
	"github.com/kode4food/ale/internal/compiler/encoder"
	"github.com/kode4food/ale/internal/debug"
	"github.com/kode4food/ale/internal/runtime/isa"
	"github.com/kode4food/ale/pkg/data"
	"github.com/kode4food/ale/pkg/env"
)

// Symbol encodes a symbol retrieval
func Symbol(e encoder.Encoder, s data.Symbol) {
	if l, ok := s.(data.Local); ok {
		resolveLocal(e, l)
		return
	}
	resolveGlobal(e, s)
}

// ReferenceSymbol encodes a potential symbol retrieval and dereference
func ReferenceSymbol(e encoder.Encoder, s data.Symbol) {
	switch s := s.(type) {
	case data.Local:
		c := resolveLocal(e, s)
		if c != nil && c.Type == encoder.ReferenceCell {
			e.Emit(isa.Deref)
		}
	default:
		resolveGlobal(e, s)
	}
}

func resolveLocal(e encoder.Encoder, l data.Local) *encoder.ScopedCell {
	if s, ok := e.ResolveScoped(l); ok {
		switch s.Scope {
		case encoder.LocalScope:
			c, _ := e.ResolveLocal(l)
			e.Emit(isa.Load, c.Index)
		case encoder.ArgScope:
			c, _ := e.ResolveParam(l)
			if c.Type == encoder.RestCell {
				e.Emit(isa.RestArg, c.Index)
			} else {
				e.Emit(isa.Arg, c.Index)
			}
		case encoder.ClosureScope:
			c, _ := e.ResolveClosure(l)
			e.Emit(isa.Closure, c.Index)
		default:
			panic(debug.ProgrammerError("unknown scope type"))
		}
		return s
	}
	resolveGlobal(e, l)
	return nil
}

func resolveGlobal(e encoder.Encoder, s data.Symbol) {
	if entry := env.MustResolveSymbol(e.Globals(), s); entry.IsBound() {
		v, _ := entry.Value()
		Literal(e, v)
		return
	}
	Literal(e, s)
	e.Emit(isa.Resolve)
}
