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
		sync.RWMutex
		environment *Environment
		domain      data.Name
		entries     entries
	}

	anonymous struct {
		Namespace
	}

	entry struct {
		sync.RWMutex
		owner Namespace
		name  data.Name
		value data.Value
		bound bool
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
	ns.RLock()
	defer ns.RUnlock()
	return ns.entries.publicNames().Sorted()
}

func (ns *namespace) Declare(n data.Name) Entry {
	ns.Lock()
	defer ns.Unlock()
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
	ns.RLock()
	defer ns.RUnlock()
	if res, ok := ns.entries[n]; ok {
		return res, true
	}
	return nil, false
}

func (ns *namespace) Snapshot(e *Environment) Namespace {
	ns.RLock()
	defer ns.RUnlock()

	res := &namespace{
		environment: e,
		domain:      ns.domain,
		entries:     make(entries, len(ns.entries)),
	}
	for k, v := range ns.entries {
		res.entries[k] = v.snapshot(res)
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

func (e *entry) snapshot(owner Namespace) *entry {
	e.RLock()
	defer e.RUnlock()
	return &entry{
		owner: owner,
		name:  e.name,
		value: e.value,
		bound: e.bound,
	}
}

func (e *entry) Owner() Namespace {
	return e.owner
}

func (e *entry) Name() data.Name {
	return e.name
}

func (e *entry) Value() data.Value {
	e.RLock()
	defer e.RUnlock()
	if e.bound {
		return e.value
	}
	panic(fmt.Errorf(ErrNameNotBound, e.name))
}

func (e *entry) IsBound() bool {
	e.RLock()
	defer e.RUnlock()
	return e.bound
}

func (e *entry) Bind(v data.Value) {
	e.Lock()
	defer e.Unlock()
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
