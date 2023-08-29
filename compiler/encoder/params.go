package encoder

import (
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/runtime/isa"
)

type paramStack []IndexedCells

func (e *encoder) PushParams(names data.Names, rest bool) {
	cells := make(IndexedCells, len(names))
	for i, n := range names {
		c := newCell(ValueCell, n)
		cells[i] = newIndexedCell(isa.Index(i), c)
	}
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

func (e *encoder) ResolveParam(n data.Name) (*IndexedCell, bool) {
	params := e.params
	for i := len(params) - 1; i >= 0; i-- {
		a := params[i]
		if c, ok := resolveParam(a, n); ok {
			return c, ok
		}
	}
	return nil, false
}

func resolveParam(cells IndexedCells, lookup data.Name) (*IndexedCell, bool) {
	for _, c := range cells {
		if c.Name == lookup {
			return c, true
		}
	}
	return nil, false
}
