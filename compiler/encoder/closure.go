package encoder

import (
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/runtime/isa"
	"github.com/kode4food/comb/slices"
)

// Closure calculates the enclosed names for this encoder
func (e *encoder) Closure() IndexedCells {
	return e.closure
}

func (e *encoder) ResolveClosure(n data.Local) (*IndexedCell, bool) {
	c, ok := slices.Find(e.closure, func(c *IndexedCell) bool {
		return c.Name == n
	})
	if ok {
		return c, ok
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
