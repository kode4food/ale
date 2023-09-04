package encoder

import (
	"errors"
	"fmt"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/runtime/isa"
)

// Locals track local variable assignments
type Locals map[data.LocalSymbol]*IndexedCell

// Error messages
const (
	ErrNoLocalScope  = "no local scopes have been pushed"
	ErrDuplicateName = "name duplicated in scope: %s"
)

func (e *encoder) LocalCount() isa.Operand {
	return e.maxLocal
}

func (e *encoder) PushLocals() {
	e.locals = append(e.locals, Locals{})
}

func (e *encoder) PopLocals() {
	if len(e.locals) == 1 {
		panic(errors.New(ErrNoLocalScope))
	}
	scope := e.peekLocals()
	e.nextLocal -= isa.Operand(len(scope))
	scopes := e.locals
	e.locals = scopes[0 : len(scopes)-1]
}

func (e *encoder) peekLocals() Locals {
	scopes := e.locals
	tailPos := len(scopes) - 1
	return scopes[tailPos]
}

func (e *encoder) allocLocal() isa.Operand {
	idx := e.nextLocal
	e.nextLocal++
	if e.nextLocal > e.maxLocal {
		e.maxLocal = e.nextLocal
	}
	return idx
}

func (e *encoder) AddLocal(n data.LocalSymbol, t CellType) *IndexedCell {
	scope := e.peekLocals()
	if _, ok := scope[n]; ok {
		panic(fmt.Errorf(ErrDuplicateName, n))
	}
	c := newCell(t, n)
	res := newIndexedCell(e.allocLocal(), c)
	scope[n] = res
	return res
}

func (e *encoder) ResolveLocal(n data.LocalSymbol) (*IndexedCell, bool) {
	scopes := e.locals
	for i := len(scopes) - 1; i >= 0; i-- {
		scope := scopes[i]
		if l, ok := scope[n]; ok {
			return l, true
		}
	}
	return nil, false
}
