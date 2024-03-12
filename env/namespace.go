package env

import (
	"sync"

	"github.com/kode4food/ale/data"
)

type (
	// Namespace represents a namespace
	Namespace interface {
		Environment() *Environment
		Domain() data.Local
		Declared() data.Locals
		Declare(data.Local) Entry
		Private(data.Local) Entry
		Resolve(data.Local) (Entry, bool)
		Snapshot(*Environment) (Namespace, error)
	}

	namespace struct {
		entries     entries
		environment *Environment
		domain      data.Local
		sync.RWMutex
	}

	anonymous struct {
		Namespace
	}
)

func (ns *namespace) Environment() *Environment {
	return ns.environment
}

func (ns *namespace) Domain() data.Local {
	return ns.domain
}

func (ns *namespace) Declared() data.Locals {
	ns.RLock()
	defer ns.RUnlock()
	var res data.Locals
	for _, e := range ns.entries {
		if !e.IsPrivate() {
			res = append(res, e.Name())
		}
	}
	return res.Sorted()
}

func (ns *namespace) Declare(n data.Local) Entry {
	return ns.declare(n)
}

func (ns *namespace) Private(n data.Local) Entry {
	e := ns.declare(n)
	e.markPrivate()
	return e
}

func (ns *namespace) declare(n data.Local) *entry {
	ns.Lock()
	defer ns.Unlock()
	if res, ok := ns.entries[n]; ok {
		return res
	}
	e := &entry{
		owner: ns,
		name:  n,
		value: data.Null,
	}
	ns.entries[n] = e
	return e
}

func (ns *namespace) Resolve(n data.Local) (Entry, bool) {
	ns.RLock()
	defer ns.RUnlock()
	if e, ok := ns.entries[n]; ok {
		e.markResolved()
		return e, true
	}
	return nil, false
}

func (ns *namespace) Snapshot(e *Environment) (Namespace, error) {
	ns.RLock()
	defer ns.RUnlock()

	res := &namespace{
		environment: e,
		domain:      ns.domain,
		entries:     make(entries, len(ns.entries)),
	}
	for k, v := range ns.entries {
		s, err := v.snapshot(res)
		if err != nil {
			return nil, err
		}
		res.entries[k] = s
	}
	return res, nil
}

func resolvePublic(from, in Namespace, n data.Local) (Entry, bool) {
	if e, ok := in.Resolve(n); ok && (from == in || !e.IsPrivate()) {
		return e, ok
	}
	return nil, false
}
