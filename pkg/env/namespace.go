package env

import (
	"fmt"
	"sync"

	"github.com/kode4food/ale/pkg/data"
)

type (
	// Namespace represents a namespace
	Namespace interface {
		Environment() *Environment
		Domain() data.Local
		Declared() data.Locals
		Declare(data.Local) Entry
		Private(data.Local) Entry
		Resolve(data.Local) (Entry, error)
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
	res := data.Locals{}
	for _, e := range ns.entries {
		if !e.IsPrivate() {
			res = append(res, e.Name())
		}
	}
	ns.RUnlock()
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
	if res, ok := ns.entries[n]; ok {
		ns.Unlock()
		return res
	}
	e := &entry{
		owner: ns,
		name:  n,
		value: data.Null,
	}
	ns.entries[n] = e
	ns.Unlock()
	return e
}

func (ns *namespace) Resolve(n data.Local) (Entry, error) {
	ns.RLock()
	if e, ok := ns.entries[n]; ok {
		e.markResolved()
		ns.RUnlock()
		return e, nil
	}
	ns.RUnlock()
	return nil, fmt.Errorf(ErrSymbolNotDeclared, n)
}

func (ns *namespace) Snapshot(e *Environment) (Namespace, error) {
	ns.RLock()
	res := &namespace{
		environment: e,
		domain:      ns.domain,
		entries:     make(entries, len(ns.entries)),
	}
	for k, v := range ns.entries {
		s, err := v.snapshot(res)
		if err != nil {
			ns.RUnlock()
			return nil, err
		}
		res.entries[k] = s
	}
	ns.RUnlock()
	return res, nil
}

func resolvePublic(from, in Namespace, n data.Local) (Entry, error) {
	if e, err := in.Resolve(n); err == nil && (from == in || !e.IsPrivate()) {
		return e, nil
	}
	return nil, fmt.Errorf(ErrSymbolNotDeclared, n)
}
