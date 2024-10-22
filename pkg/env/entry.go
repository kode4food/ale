package env

import (
	"fmt"
	"sync"

	"github.com/kode4food/ale/pkg/data"
)

type (
	// Entry represents a namespace entry
	Entry interface {
		Owner() Namespace
		Name() data.Local
		Value() (data.Value, error)
		Bind(data.Value) error
		IsBound() bool
		IsPrivate() bool
	}

	entry struct {
		owner Namespace
		value data.Value
		name  data.Local
		flags entryFlag
		sync.RWMutex
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
	if e.hasFlag(resolved) && !e.hasFlag(bound) {
		e.RUnlock()
		return nil, fmt.Errorf(ErrSnapshotIncomplete, e.name)
	}
	res := &entry{
		owner: owner,
		name:  e.name,
		value: e.value,
		flags: e.flags,
	}
	e.RUnlock()
	return res, nil
}

func (e *entry) markResolved() {
	e.Lock()
	e.setFlag(resolved)
	e.Unlock()
}

func (e *entry) markPrivate() {
	e.Lock()
	e.setFlag(private)
	e.Unlock()
}

func (e *entry) Owner() Namespace {
	return e.owner
}

func (e *entry) Name() data.Local {
	return e.name
}

func (e *entry) Value() (data.Value, error) {
	e.RLock()
	if e.hasFlag(bound) {
		res := e.value
		e.RUnlock()
		return res, nil
	}
	e.RUnlock()
	return nil, fmt.Errorf(ErrNameNotBound, e.name)
}

func (e *entry) Bind(v data.Value) error {
	e.Lock()
	if e.hasFlag(bound) {
		e.Unlock()
		return fmt.Errorf(ErrNameAlreadyBound, e.name)
	}
	e.value = v
	e.setFlag(bound)
	e.Unlock()
	return nil
}

func (e *entry) IsBound() bool {
	e.RLock()
	res := e.hasFlag(bound)
	e.RUnlock()
	return res
}

func (e *entry) IsPrivate() bool {
	e.RLock()
	res := e.hasFlag(private)
	e.RUnlock()
	return res
}

func (e *entry) hasFlag(flag entryFlag) bool {
	return e.flags&flag != 0
}

func (e *entry) setFlag(flag entryFlag) {
	e.flags |= flag
}
