package encoder

import (
	"errors"
	"fmt"

	"github.com/kode4food/ale/internal/runtime/isa"
	"github.com/kode4food/ale/pkg/data"
)

// Locals track local variable assignments
type Locals map[data.Local]*IndexedCell

const (
	// ErrNoLocalScope is raised when no local scope has been created
	ErrNoLocalScope = "no local scopes have been pushed"

	// ErrDuplicateName is raised when an attempt is made to register a
	// duplicated name within the same local scope
	ErrDuplicateName = "name duplicated in scope: %s"
)

func (e *encoder) PushLocals() {
	e.locals = append(e.locals, Locals{})
}

func (e *encoder) PopLocals() error {
	if len(e.locals) == 1 {
		return errors.New(ErrNoLocalScope)
	}
	scope := e.peekLocals()
	e.nextLocal -= isa.Operand(len(scope))
	scopes := e.locals
	e.locals = scopes[0 : len(scopes)-1]
	return nil
}

func (e *encoder) peekLocals() Locals {
	scopes := e.locals
	tailPos := len(scopes) - 1
	return scopes[tailPos]
}

func (e *encoder) allocLocal() isa.Operand {
	idx := e.nextLocal
	e.nextLocal++
	return idx
}

func (e *encoder) AddLocal(n data.Local, t CellType) (*IndexedCell, error) {
	scope := e.peekLocals()
	if _, ok := scope[n]; ok {
		return nil, fmt.Errorf(ErrDuplicateName, n)
	}
	c := newCell(t, n)
	res := newIndexedCell(e.allocLocal(), c)
	scope[n] = res
	return res, nil
}

func (e *encoder) ResolveLocal(n data.Local) (*IndexedCell, bool) {
	scopes := e.locals
	for i := len(scopes) - 1; i >= 0; i-- {
		scope := scopes[i]
		if l, ok := scope[n]; ok {
			return l, true
		}
	}
	return nil, false
}
