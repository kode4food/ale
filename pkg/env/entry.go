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
		sync.RWMutex
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
)

const (
	public  entryFlag = 0
	private entryFlag = 1 << iota
	bound
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
	return e.hasFlag(bound)
}

func (e *entry) IsPrivate() bool {
	return e.hasFlag(private)
}

func (e *entry) snapshot() *entry {
	if e.hasFlag(bound) {
		return e
	}

	return &entry{
		name:  e.name,
		value: e.value,
		flags: e.flags,
	}
}

func (e *entry) hasFlag(flag entryFlag) bool {
	return flag == 0 || atomic.LoadUint64((*uint64)(&e.flags))&uint64(flag) != 0
}

func (e *entry) setFlag(flag entryFlag) {
	atomic.OrUint64((*uint64)(&e.flags), uint64(flag))
}
