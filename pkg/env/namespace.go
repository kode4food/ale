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
	defer ns.RUnlock()
	res := data.Locals{}
	for _, e := range ns.entries {
		if !e.IsPrivate() {
			res = append(res, e.Name())
		}
	}
	return res.Sorted()
}

func (ns *namespace) Public(n data.Local) (Entry, error) {
	return ns.declare(n, public)
}

func (ns *namespace) Private(n data.Local) (Entry, error) {
	return ns.declare(n, private)
}

func (ns *namespace) declare(n data.Local, f entryFlag) (*entry, error) {
	ns.Lock()
	defer ns.Unlock()
	if e, ok := ns.entries[n]; ok {
		if e.hasFlag(f) {
			return e, nil
		}
		return nil, fmt.Errorf(ErrNameAlreadyDeclared, n)
	}
	e := &entry{
		name:  n,
		flags: f,
	}
	ns.entries[n] = e
	return e, nil
}

func (ns *namespace) Resolve(n data.Local) (Entry, Namespace, error) {
	ns.RLock()
	defer ns.RUnlock()
	if e, ok := ns.entries[n]; ok {
		e.markResolved()
		return e, ns, nil
	}
	return nil, nil, fmt.Errorf(ErrNameNotDeclared, n)
}

func (ns *namespace) Snapshot(e *Environment) (Namespace, error) {
	ns.RLock()
	defer ns.RUnlock()
	res := &namespace{
		environment: e,
		domain:      ns.domain,
		entries:     make(entries, len(ns.entries)),
	}
	var err error
	for k, v := range ns.entries {
		if res.entries[k], err = v.snapshot(); err != nil {
			return nil, err
		}
	}
	return res, nil
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
