package encoder

import "gitlab.com/kode4food/ale/api"

// Scope describes the scope of a name
type Scope int

//go:generate stringer -type=Scope -linecomment
const (
	LocalScope Scope = iota
	ArgScope
	NameScope
	ClosureScope
)

func (e *encoder) ResolveScope(l api.LocalSymbol) (Scope, bool) {
	if _, ok := e.ResolveLocal(l); ok {
		return LocalScope, true
	}
	if _, _, ok := e.ResolveArg(l); ok {
		return ArgScope, true
	}
	if e.Name() == l.Name() {
		return NameScope, true
	}
	if e.parent != nil {
		if _, ok := e.parent.ResolveScope(l); ok {
			return ClosureScope, true
		}
	}
	return -1, false
}

func (e *encoder) InScope(l api.LocalSymbol) bool {
	_, ok := e.ResolveScope(l)
	return ok
}
