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
		flags entryFlags
	}

	entries map[data.Local]*entry

	entryFlags uint32
)

const (
	resolved entryFlags = 1 << iota
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

func (e *entry) hasFlag(flag entryFlags) bool {
	return e.flags&flag != 0
}

func (e *entry) setFlag(flag entryFlags) {
	e.flags |= flag
}
