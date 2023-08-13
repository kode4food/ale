package env

import (
	"fmt"
	"sync"

	"github.com/kode4food/ale/data"
)

type (
	// Namespace represents a namespace
	Namespace interface {
		Environment() *Environment
		Domain() data.Name
		Declared() data.Names
		Declare(data.Name) Entry
		Resolve(data.Name) (Entry, bool)
		Snapshot(*Environment) Namespace
	}

	// Entry represents a namespace entry
	Entry interface {
		Owner() Namespace
		Name() data.Name
		Value() data.Value
		IsBound() bool
		Bind(data.Value)
	}

	namespace struct {
		environment *Environment
		domain      data.Name
		entries     entries
		mutex       sync.RWMutex
	}

	anonymous struct {
		Namespace
	}

	entry struct {
		owner Namespace
		name  data.Name
		value data.Value
		bound bool
		mutex sync.RWMutex
	}

	entries map[data.Name]*entry
)

// Error messages
const (
	ErrNameAlreadyBound = "name is already bound in namespace: %s"
	ErrNameNotBound     = "name is not bound in namespace: %s"
)

func (ns *namespace) Environment() *Environment {
	return ns.environment
}

func (ns *namespace) Domain() data.Name {
	return ns.domain
}

func (ns *namespace) Declared() data.Names {
	ns.mutex.RLock()
	defer ns.mutex.RUnlock()
	return ns.entries.publicNames().Sorted()
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

func (ns *namespace) Snapshot(e *Environment) Namespace {
	ns.mutex.RLock()
	defer ns.mutex.RUnlock()

	res := &namespace{
		environment: e,
		domain:      ns.domain,
		entries:     make(entries, len(ns.entries)),
	}
	for k, v := range ns.entries {
		res.entries[k] = &entry{
			owner: res,
			name:  v.name,
			value: v.value,
			bound: v.bound,
		}
	}
	return res
}

func (e entries) publicNames() data.Names {
	res := make(data.Names, 0, len(e))
	for k := range e {
		if !isPrivateName(k) {
			res = append(res, k)
		}
	}
	return res
}

func (e *entry) Owner() Namespace {
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
	panic(fmt.Errorf(ErrNameNotBound, e.name))
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
		panic(fmt.Errorf(ErrNameAlreadyBound, e.name))
	}
	e.value = v
	e.bound = true
}

func resolvePublic(from, in Namespace, n data.Name) (Entry, bool) {
	if isPrivateName(n) && from != in {
		return nil, false
	}
	return in.Resolve(n)
}

func isPrivateName(n data.Name) bool {
	return len(n) > 1 && n[0] == '^'
}
