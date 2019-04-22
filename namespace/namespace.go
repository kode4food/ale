package namespace

import (
	"fmt"
	"sync"

	"gitlab.com/kode4food/ale/data"
)

type (
	// Type represents a namespace
	Type interface {
		Manager() *Manager
		Domain() data.Name
		Declare(data.Name) Entry
		Resolve(data.Name) (Entry, bool)
	}

	// Entry represents a namespace entry
	Entry interface {
		Owner() Type
		Name() data.Name
		Value() data.Value
		IsBound() bool
		Bind(data.Value)
	}

	namespace struct {
		manager *Manager
		domain  data.Name
		entries entries
		mutex   sync.RWMutex
	}

	anonymous struct {
		Type
	}

	entry struct {
		owner Type
		name  data.Name
		value data.Value
		bound bool
		mutex sync.RWMutex
	}

	entries map[data.Name]*entry
)

// Error messages
const (
	NameAlreadyBound = "name is already bound in namespace: %s"
	NameNotBound     = "name is not bound in namespace: %s"
)

func (ns *namespace) Manager() *Manager {
	return ns.manager
}

func (ns *namespace) Domain() data.Name {
	return ns.domain
}

func (ns *namespace) Declare(n data.Name) Entry {
	ns.mutex.Lock()
	defer ns.mutex.Unlock()
	if res, ok := ns.entries[n]; ok {
		return res
	}
	e := &entry{
		owner: ns,
		name:  n,
		value: data.Nil,
		bound: false,
	}
	ns.entries[n] = e
	return e
}

func (ns *namespace) Resolve(n data.Name) (Entry, bool) {
	ns.mutex.RLock()
	defer ns.mutex.RUnlock()
	if res, ok := ns.entries[n]; ok {
		return res, true
	}
	return nil, false
}

func (e *entry) Owner() Type {
	return e.owner
}

func (e *entry) Name() data.Name {
	return e.name
}

func (e *entry) Value() data.Value {
	e.mutex.RLock()
	defer e.mutex.RUnlock()
	if e.bound {
		return e.value
	}
	panic(fmt.Errorf(NameNotBound, e.name))
}

func (e *entry) IsBound() bool {
	e.mutex.RLock()
	defer e.mutex.RUnlock()
	return e.bound
}

func (e *entry) Bind(v data.Value) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	if e.bound {
		panic(fmt.Errorf(NameAlreadyBound, e.name))
	}
	e.value = v
	e.bound = true
}
