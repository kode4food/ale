package encoder

import "github.com/kode4food/ale/pkg/data"

// Scope describes the scope of a name
type Scope int

// Scope locations
const (
	LocalScope Scope = iota
	ArgScope
	ClosureScope
)

func (e *encoder) ResolveScoped(n data.Local) (*ScopedCell, bool) {
	if i, ok := e.ResolveLocal(n); ok {
		return newScopedCell(LocalScope, i.Cell), true
	}
	if _, ok := e.ResolveParam(n); ok {
		return newScopedCell(ArgScope, newCell(ValueCell, n)), true
	}
	if e.parent == nil {
		return nil, false
	}
	if s, ok := e.parent.ResolveScoped(n); ok {
		return newScopedCell(ClosureScope, s.Cell), true
	}
	return nil, false
}
