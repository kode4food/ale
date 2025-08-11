package env

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/kode4food/ale"
	"github.com/kode4food/ale/data"
)

type (
	// Entry represents a namespace entry
	Entry struct {
		name    data.Local
		private bool
		*binding
	}

	binding struct {
		value ale.Value
		bound atomic.Bool
		sync.Mutex
	}
)

const (
	// ErrNameAlreadyBound is raised when an attempt is made to bind a
	// Namespace entry that has already been bound
	ErrNameAlreadyBound = "name is already bound in namespace: %s"

	// ErrNameNotBound is raised when an attempt is made to retrieve a value
	// from a Namespace that hasn't been bound
	ErrNameNotBound = "name is not bound in namespace: %s"
)

func (e *Entry) Name() data.Local {
	return e.name
}

func (e *Entry) Value() (ale.Value, error) {
	if e.bound.Load() {
		return e.value, nil
	}
	return nil, fmt.Errorf(ErrNameNotBound, e.name)
}

func (e *Entry) Bind(v ale.Value) error {
	e.Lock()
	defer e.Unlock()
	if e.bound.Load() {
		return fmt.Errorf(ErrNameAlreadyBound, e.name)
	}
	e.value = v
	e.bound.Store(true)
	return nil
}

func (e *Entry) IsBound() bool {
	return e.bound.Load()
}

func (e *Entry) IsPrivate() bool {
	return e.private
}

func (e *Entry) snapshot() *Entry {
	if e.bound.Load() {
		return e
	}

	return &Entry{
		name:    e.name,
		private: e.private,
		binding: new(binding),
	}
}
