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
		Public(data.Local) (Entry, error)
		Private(data.Local) (Entry, error)
		Resolve(data.Local) (Entry, Namespace, error)
		Snapshot(*Environment) Namespace
	}

	namespace struct {
		entries     entries
		environment *Environment
		domain      data.Local
		sync.RWMutex
	}
)

const (
	// ErrNameAlreadyDeclared is raised when an attempt to declare a name is
	// performed that has already been declared with different privacy
	ErrNameAlreadyDeclared = "name already declared in namespace: %s"

	// ErrNameNotDeclared is raised when an attempt to forcefully resolve an
	// undeclared name in the Namespace fails
	ErrNameNotDeclared = "name not declared in namespace: %s"
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

func (ns *namespace) Public(n data.Local) (Entry, error) {
	return ns.declare(n, false)
}

func (ns *namespace) Private(n data.Local) (Entry, error) {
	return ns.declare(n, true)
}

func (ns *namespace) declare(n data.Local, makePrivate bool) (*entry, error) {
	ns.Lock()
	if e, ok := ns.entries[n]; ok {
		if e.hasFlag(private) == makePrivate {
			ns.Unlock()
			return e, nil
		}
		ns.Unlock()
		return nil, fmt.Errorf(ErrNameAlreadyDeclared, n)
	}
	e := &entry{name: n}
	if makePrivate {
		e.flags = private
	}
	ns.entries[n] = e
	ns.Unlock()
	return e, nil
}

func (ns *namespace) Resolve(n data.Local) (Entry, Namespace, error) {
	if e, ok := ns.resolve(n); ok {
		return e, ns, nil
	}
	return nil, nil, fmt.Errorf(ErrNameNotDeclared, n)
}

func (ns *namespace) resolve(n data.Local) (*entry, bool) {
	ns.RLock()
	e, ok := ns.entries[n]
	ns.RUnlock()
	return e, ok
}

func (ns *namespace) Snapshot(e *Environment) Namespace {
	ns.RLock()
	res := &namespace{
		environment: e,
		domain:      ns.domain,
		entries:     make(entries, len(ns.entries)),
	}
	for k, v := range ns.entries {
		res.entries[k] = v.snapshot()
	}
	ns.RUnlock()
	return res
}

func resolvePublic(from, in Namespace, n data.Local) (Entry, Namespace, error) {
	if e, ns, err := in.Resolve(n); err == nil {
		if !e.IsPrivate() || from == in && in == ns {
			return e, ns, nil
		}
	}
	return nil, nil, fmt.Errorf(ErrNameNotDeclared, n)
}

func BindPublic(ns Namespace, n data.Local, v data.Value) error {
	e, err := ns.Public(n)
	if err != nil {
		return err
	}
	return e.Bind(v)
}

func BindPrivate(ns Namespace, n data.Local, v data.Value) error {
	e, err := ns.Private(n)
	if err != nil {
		return err
	}
	return e.Bind(v)
}
