package encoder

import (
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/runtime/isa"
	"github.com/kode4food/comb/slices"
)

type paramStack []IndexedCells

func (e *encoder) PushParams(names data.Locals, rest bool) {
	cells := slices.IndexedMap(names, func(n data.Local, i int) *IndexedCell {
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
	al := len(params)
	e.params = params[0 : al-1]
}

func (e *encoder) ResolveParam(n data.Local) (*IndexedCell, bool) {
	params := e.params
	for i := len(params) - 1; i >= 0; i-- {
		a := params[i]
		if c, ok := resolveParam(a, n); ok {
			return c, ok
		}
	}
	return nil, false
}

func resolveParam(cells IndexedCells, lookup data.Local) (*IndexedCell, bool) {
	return slices.Find(cells, func(c *IndexedCell) bool {
		return c.Name == lookup
	})
}
