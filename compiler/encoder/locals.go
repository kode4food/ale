package encoder

import (
	"fmt"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/runtime/isa"
)

// Locals track local variable assignments
type Locals map[data.Name]*IndexedCell

// Error messages
const (
	ErrDuplicateName = "name duplicated in scope: %s"
)

func (e *encoder) LocalCount() int {
	return e.maxLocal
}

func (e *encoder) PushLocals() {
	e.locals = append(e.locals, Locals{})
}

func (e *encoder) PopLocals() {
	scope := e.peekLocals()
	e.nextLocal -= len(scope)
	scopes := e.locals
	e.locals = scopes[0 : len(scopes)-1]
}

func (e *encoder) peekLocals() Locals {
	scopes := e.locals
	tailPos := len(scopes) - 1
	return scopes[tailPos]
}

func (e *encoder) allocLocal() isa.Index {
	idx := isa.Index(e.nextLocal)
	e.nextLocal++
	if e.nextLocal > e.maxLocal {
		e.maxLocal = e.nextLocal
	}
	return idx
}

func (e *encoder) AddLocal(n data.Name, t CellType) *IndexedCell {
	scope := e.peekLocals()
	if _, ok := scope[n]; ok {
		panic(fmt.Errorf(ErrDuplicateName, n))
	}
	c := newCell(t, n)
	res := newIndexedCell(e.allocLocal(), c)
	scope[n] = res
	return res
}

func (e *encoder) ResolveLocal(n data.Name) (*IndexedCell, bool) {
	scopes := e.locals
	for i := len(scopes) - 1; i >= 0; i-- {
		scope := scopes[i]
		if l, ok := scope[n]; ok {
			return l, true
		}
	}
	return nil, false
}
