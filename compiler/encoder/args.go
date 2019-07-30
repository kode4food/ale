package encoder

import (
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/runtime/isa"
)

type argsStack []IndexedCells

func (e *encoder) PushArgs(names data.Names, rest bool) {
	cells := make(IndexedCells, len(names))
	for i, n := range names {
		c := newCell(ValueCell, n)
		cells[i] = newIndexedCell(isa.Index(i), c)
	}
	if rest {
		cells[len(cells)-1].Type = RestCell
	}
	e.args = append(e.args, cells)
}

func (e *encoder) PopArgs() {
	args := e.args
	al := len(args)
	e.args = args[0 : al-1]
}

func (e *encoder) ResolveArg(n data.Name) (*IndexedCell, bool) {
	args := e.args
	for i := len(args) - 1; i >= 0; i-- {
		a := args[i]
		if c, ok := resolveArg(a, n); ok {
			return c, ok
		}
	}
	return nil, false
}

func resolveArg(cells IndexedCells, lookup data.Name) (*IndexedCell, bool) {
	for _, c := range cells {
		if c.Name == lookup {
			return c, true
		}
	}
	return nil, false
}
