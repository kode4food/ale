package encoder

import (
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/runtime/isa"
)

// Closure calculates the enclosed names for this encoder
func (e *encoder) Closure() IndexedCells {
	return e.closure
}

func (e *encoder) ResolveClosure(n data.Local) (*IndexedCell, bool) {
	for _, c := range e.closure {
		if c.Name == n {
			return c, true
		}
	}
	return e.resolveClosureParent(n)
}

func (e *encoder) resolveClosureParent(n data.Local) (*IndexedCell, bool) {
	parent := e.parent
	if parent == nil {
		return nil, false
	}
	s, ok := parent.ResolveScoped(n)
	if !ok {
		return nil, false
	}
	closure := e.closure
	idx := isa.Operand(len(closure))
	res := newIndexedCell(idx, s.Cell)
	e.closure = append(closure, res)
	return res, true
}
