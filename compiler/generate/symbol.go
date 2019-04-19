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
	GlobalNotDeclared = "symbol not declared in namespace: %s"
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
			panic(fmt.Sprintf("unknown scope type: %s", scope))
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
	if v, ok := globals.Resolve(name); ok {
		Literal(e, v)
		return
	}
	if !globals.IsDeclared(name) {
		panic(fmt.Errorf(GlobalNotDeclared, name))
	}
	resolveFromNamespace(e, globals, l)
}

func resolveFromNamespace(e encoder.Type, ns namespace.Type, s data.Symbol) {
	name := s.Name()
	if ns.IsBound(name) {
		v, _ := ns.Resolve(name)
		Literal(e, v)
		return
	}
	Literal(e, s)
	e.Emit(isa.Resolve)
}
