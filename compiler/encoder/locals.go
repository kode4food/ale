package encoder

import (
	"fmt"

	"gitlab.com/kode4food/ale/data"
	"gitlab.com/kode4food/ale/runtime/isa"
)

// Locals tracks local variable assignments
type Locals map[data.Name]isa.Index

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

func (e *encoder) AddLocal(n data.Name) isa.Index {
	scope := e.peekLocals()
	if _, ok := scope[n]; ok {
		panic(fmt.Sprintf("name duplicated in scope: %s", n))
	}
	idx := e.allocLocal()
	scope[n] = idx
	return idx
}

func (e *encoder) ResolveLocal(l data.LocalSymbol) (isa.Index, bool) {
	n := l.Name()
	scopes := e.locals
	for i := len(scopes) - 1; i >= 0; i-- {
		scope := scopes[i]
		if i, ok := scope[n]; ok {
			return i, true
		}
	}
	return 0, false
}
