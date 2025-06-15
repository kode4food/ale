package generate

import (
	"github.com/kode4food/ale/internal/compiler/encoder"
	"github.com/kode4food/ale/internal/debug"
	"github.com/kode4food/ale/internal/runtime/isa"
	"github.com/kode4food/ale/pkg/data"
	"github.com/kode4food/ale/pkg/env"
)

// Global encodes a global symbol constant or retrieval, depending on whether
// the symbol is already bound in the environment
func Global(e encoder.Encoder, s data.Symbol) error {
	entry, _, err := env.ResolveSymbol(e.Globals(), s)
	if err != nil {
		return err
	}
	if entry.IsBound() {
		v, _ := entry.Value()
		return Literal(e, v)
	}
	if err := Literal(e, s); err != nil {
		return err
	}
	e.Emit(isa.EnvValue)
	return nil
}

// Reference encodes a potential retrieval and dereference
func Reference(e encoder.Encoder, l data.Local) error {
	c, err := resolveLocal(e, l)
	if err != nil {
		return err
	}
	if c != nil && c.Type == encoder.ReferenceCell {
		e.Emit(isa.RefValue)
	}
	return nil
}

// Local encodes a local retrieval, but not dereference
func Local(e encoder.Encoder, l data.Local) error {
	_, err := resolveLocal(e, l)
	return err
}

func resolveLocal(
	e encoder.Encoder, l data.Local,
) (*encoder.ScopedCell, error) {
	if s, ok := e.ResolveScoped(l); ok {
		switch s.Scope {
		case encoder.LocalScope:
			c, _ := e.ResolveLocal(l)
			e.Emit(isa.Load, c.Index)
		case encoder.ArgScope:
			c, _ := e.ResolveParam(l)
			if c.Type == encoder.RestCell {
				e.Emit(isa.ArgsRest, c.Index)
			} else {
				e.Emit(isa.Arg, c.Index)
			}
		case encoder.ClosureScope:
			c, _ := e.ResolveClosure(l)
			e.Emit(isa.Closure, c.Index)
		default:
			panic(debug.ProgrammerError("unknown scope type"))
		}
		return s, nil
	}
	return nil, Global(e, l)
}
