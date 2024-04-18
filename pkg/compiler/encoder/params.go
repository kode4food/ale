package encoder

import (
	"github.com/kode4food/ale/pkg/data"
	"github.com/kode4food/ale/pkg/runtime/isa"
	"github.com/kode4food/comb/basics"
)

type paramStack []IndexedCells

func (e *encoder) PushParams(names data.Locals, rest bool) {
	cells := basics.IndexedMap(names, func(n data.Local, i int) *IndexedCell {
		c := newCell(ValueCell, n)
		return newIndexedCell(isa.Operand(i), c)
	})
	if rest {
		cells[len(cells)-1].Type = RestCell
	}
	e.params = append(e.params, cells)
}

func (e *encoder) PopParams() {
	params := e.params
	pl := len(params)
	e.params = params[0 : pl-1]
}

func (e *encoder) ResolveParam(n data.Local) (*IndexedCell, bool) {
	params := e.params
	for i := len(params) - 1; i >= 0; i-- {
		p := params[i]
		if c, ok := resolveParam(p, n); ok {
			return c, ok
		}
	}
	return nil, false
}

func resolveParam(cells IndexedCells, lookup data.Local) (*IndexedCell, bool) {
	return basics.Find(cells, func(c *IndexedCell) bool {
		return c.Name == lookup
	})
}
