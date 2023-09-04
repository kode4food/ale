package env

import (
	"sync"

	"github.com/kode4food/ale/data"
)

type (
	// Namespace represents a namespace
	Namespace interface {
		Environment() *Environment
		Domain() data.LocalSymbol
		Declared() data.LocalSymbols
		Declare(data.LocalSymbol) Entry
		Private(data.LocalSymbol) Entry
		Resolve(data.LocalSymbol) (Entry, bool)
		Snapshot(*Environment) (Namespace, error)
	}

	namespace struct {
		sync.RWMutex
		environment *Environment
		domain      data.LocalSymbol
		entries     entries
	}

	anonymous struct {
		Namespace
	}
)

// Error messages
const (
	ErrNameAlreadyBound   = "name is already bound in namespace: %s"
	ErrNameNotBound       = "name is not bound in namespace: %s"
	ErrSnapshotIncomplete = "can't snapshot environment. entry not bound: %s"
)

func (ns *namespace) Environment() *Environment {
	return ns.environment
}

func (ns *namespace) Domain() data.LocalSymbol {
	return ns.domain
}

func (ns *namespace) Declared() data.LocalSymbols {
	ns.RLock()
	defer ns.RUnlock()
	var res data.LocalSymbols
	for _, e := range ns.entries {
		if !e.IsPrivate() {
			res = append(res, e.Name())
		}
	}
	return res.Sorted()
}

func (ns *namespace) Declare(n data.LocalSymbol) Entry {
	return ns.declare(n)
}

func (ns *namespace) Private(n data.LocalSymbol) Entry {
	e := ns.declare(n)
	e.markPrivate()
	return e
}

func (ns *namespace) declare(n data.LocalSymbol) *entry {
	ns.Lock()
	defer ns.Unlock()
	if res, ok := ns.entries[n]; ok {
		return res
	}
	e := &entry{
		owner: ns,
		name:  n,
		value: data.Nil,
	}
	ns.entries[n] = e
	return e
}

func (ns *namespace) Resolve(n data.LocalSymbol) (Entry, bool) {
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

func resolvePublic(from, in Namespace, n data.LocalSymbol) (Entry, bool) {
	if e, ok := in.Resolve(n); ok && (from == in || !e.IsPrivate()) {
		return e, ok
	}
	return nil, false
}
