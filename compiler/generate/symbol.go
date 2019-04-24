package generate

import (
	"fmt"

	"gitlab.com/kode4food/ale/compiler/encoder"
	"gitlab.com/kode4food/ale/data"
	"gitlab.com/kode4food/ale/namespace"
	"gitlab.com/kode4food/ale/runtime/isa"
)

// Error messages
const (
	SymbolNotDeclared = "symbol not declared in namespace: %s"
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
	if l, ok := s.(data.LocalSymbol); ok {
		resolveFromEncoder(e, l)
		return
	}
	q := s.(data.QualifiedSymbol)
	manager := e.Globals().Manager()
	ns := manager.GetQualified(q.Domain())
	resolveFromNamespace(e, ns, q)
}

func resolveFromEncoder(e encoder.Type, l data.LocalSymbol) {
	globals := e.Globals()
	name := l.Name()
	ge, ok := globals.Resolve(name)
	if !ok {
		panic(fmt.Errorf(SymbolNotDeclared, name))
	}
	if ge.IsBound() {
		Literal(e, ge.Value())
		return
	}
	resolveFromNamespace(e, globals, l)
}

func resolveFromNamespace(e encoder.Type, ns namespace.Type, s data.Symbol) {
	name := s.Name()
	if ne, ok := ns.Resolve(name); ok && ne.Owner() == ns && ne.IsBound() {
		Literal(e, ne.Value())
		return
	}
	Literal(e, s)
	e.Emit(isa.Resolve)
}
