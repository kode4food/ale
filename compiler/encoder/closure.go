package encoder

import (
	"gitlab.com/kode4food/ale/data"
	"gitlab.com/kode4food/ale/runtime/isa"
)

// Closure calculates the enclosed names for this encoder
func (e *encoder) Closure() data.Names {
	return e.closure
}

func (e *encoder) ResolveClosure(l data.LocalSymbol) (isa.Index, bool) {
	closure := e.closure
	lookup := l.Name()
	for idx, n := range closure {
		if n == lookup {
			return isa.Index(idx), true
		}
	}
	parent := e.parent
	if parent != nil && parent.InScope(l) {
		res := len(closure)
		e.closure = append(closure, lookup)
		return isa.Index(res), true
	}
	return 0, false
}
