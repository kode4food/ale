package env

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/kode4food/ale/pkg/data"
)

type (
	// Entry represents a namespace entry
	Entry interface {
		Name() data.Local
		Value() (data.Value, error)
		Bind(data.Value) error
		IsBound() bool
		IsPrivate() bool
	}

	entry struct {
		value data.Value
		name  data.Local
		flags entryFlag
		sync.Mutex
	}

	entries map[data.Local]*entry

	entryFlag uint64
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
	public entryFlag = 0

	resolved entryFlag = 1 << iota
	bound
	private
)

func (e *entry) Name() data.Local {
	return e.name
}

func (e *entry) Value() (data.Value, error) {
	if e.hasFlag(bound) {
		return e.value, nil
	}
	return nil, fmt.Errorf(ErrNameNotBound, e.name)
}

func (e *entry) Bind(v data.Value) error {
	e.Lock()
	defer e.Unlock()
	if e.hasFlag(bound) {
		return fmt.Errorf(ErrNameAlreadyBound, e.name)
	}
	e.value = v
	e.setFlag(bound | resolved)
	return nil
}

func (e *entry) IsBound() bool {
	return e.hasFlag(bound)
}

func (e *entry) IsPrivate() bool {
	return e.hasFlag(private)
}

func (e *entry) snapshot() (*entry, error) {
	e.Lock()
	defer e.Unlock()
	if e.hasFlag(bound) {
		return e, nil
	}

	if e.hasFlag(resolved) {
		return nil, fmt.Errorf(ErrSnapshotIncomplete, e.name)
	}

	return &entry{
		name:  e.name,
		value: e.value,
		flags: e.flags,
	}, nil
}

func (e *entry) markResolved() {
	if !e.hasFlag(resolved) {
		e.setFlag(resolved)
	}
}

func (e *entry) hasFlag(flag entryFlag) bool {
	return flag == 0 || atomic.LoadUint64((*uint64)(&e.flags))&uint64(flag) != 0
}

func (e *entry) setFlag(flag entryFlag) {
	atomic.OrUint64((*uint64)(&e.flags), uint64(flag))
}
