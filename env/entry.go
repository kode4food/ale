package env

import (
	"fmt"
	"sync"

	"github.com/kode4food/ale/data"
)

type (
	// Entry represents a namespace entry
	Entry interface {
		Owner() Namespace
		Name() data.Local
		Value() data.Value
		Bind(data.Value)
		IsBound() bool
		IsPrivate() bool
	}

	entry struct {
		sync.RWMutex
		owner Namespace
		name  data.Local
		value data.Value
		flags entryFlag
	}

	entries map[data.Local]*entry

	entryFlag uint
)

const (
	// ErrNameAlreadyBound is raised when an attempt is made to bind a
	// Namespace entry that has already been bound
	ErrNameAlreadyBound = "name is already bound in namespace: %s"

	// ErrNameNotBound is raised when an attempt is mode to retrieve a value
	// from a Namespace that hasn't been bound
	ErrNameNotBound = "name is not bound in namespace: %s"

	// ErrSnapshotIncomplete is raised when an attempt is made to create a
	// Namespace snapshot in a situation where an unbound entry has been
	// retrieved
	ErrSnapshotIncomplete = "can't snapshot environment. entry not bound: %s"
)

const (
	resolved entryFlag = 1 << iota
	bound
	private
)

func (e *entry) snapshot(owner Namespace) (*entry, error) {
	e.RLock()
	defer e.RUnlock()
	if e.hasFlag(resolved) && !e.hasFlag(bound) {
		return nil, fmt.Errorf(ErrSnapshotIncomplete, e.name)
	}
	return &entry{
		owner: owner,
		name:  e.name,
		value: e.value,
		flags: e.flags,
	}, nil
}

func (e *entry) markResolved() {
	e.Lock()
	defer e.Unlock()
	e.setFlag(resolved)
}

func (e *entry) markPrivate() {
	e.Lock()
	defer e.Unlock()
	e.setFlag(private)
}

func (e *entry) Owner() Namespace {
	return e.owner
}

func (e *entry) Name() data.Local {
	return e.name
}

func (e *entry) Value() data.Value {
	e.RLock()
	defer e.RUnlock()
	if e.hasFlag(bound) {
		return e.value
	}
	panic(fmt.Errorf(ErrNameNotBound, e.name))
}

func (e *entry) Bind(v data.Value) {
	e.Lock()
	defer e.Unlock()
	if e.hasFlag(bound) {
		panic(fmt.Errorf(ErrNameAlreadyBound, e.name))
	}
	e.value = v
	e.setFlag(bound)
}

func (e *entry) IsBound() bool {
	e.RLock()
	defer e.RUnlock()
	return e.hasFlag(bound)
}

func (e *entry) IsPrivate() bool {
	e.RLock()
	defer e.RUnlock()
	return e.hasFlag(private)
}

func (e *entry) hasFlag(flag entryFlag) bool {
	return e.flags&flag != 0
}

func (e *entry) setFlag(flag entryFlag) {
	e.flags |= flag
}
