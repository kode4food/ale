package env

import (
	"fmt"
	"sync"

	"github.com/kode4food/ale/internal/basics"
	"github.com/kode4food/ale/pkg/data"
)

type (
	// Namespace represents a namespace
	Namespace interface {
		Environment() *Environment
		Domain() data.Local
		Declared() data.Locals
		Public(data.Local) (*Entry, error)
		Private(data.Local) (*Entry, error)
		Resolve(data.Local) (*Entry, Namespace, error)
		Snapshot(*Environment) Namespace
		Import(map[data.Local]*Entry) error
	}

	Binder func(ns Namespace, n data.Local, v data.Value) error

	namespace struct {
		entries     entries
		environment *Environment
		domain      data.Local
		sync.RWMutex
	}

	entries map[data.Local]*Entry
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
	res := make(data.Locals, 0, len(ns.entries))
	for _, e := range ns.entries {
		if e.IsPrivate() {
			continue
		}
		res = append(res, e.Name())
	}
	return res
}

func (ns *namespace) Public(n data.Local) (*Entry, error) {
	return ns.declare(n, false)
}

func (ns *namespace) Private(n data.Local) (*Entry, error) {
	return ns.declare(n, true)
}

func (ns *namespace) declare(n data.Local, asPrivate bool) (*Entry, error) {
	ns.Lock()
	defer ns.Unlock()
	if e, ok := ns.entries[n]; ok {
		if e.private == asPrivate {
			return e, nil
		}
		return nil, fmt.Errorf(ErrNameAlreadyDeclared, n)
	}
	e := &Entry{
		name:    n,
		private: asPrivate,
		binding: &binding{},
	}
	ns.entries[n] = e
	return e, nil
}

func (ns *namespace) Resolve(n data.Local) (*Entry, Namespace, error) {
	if e, ok := ns.resolve(n); ok {
		return e, ns, nil
	}
	return nil, nil, fmt.Errorf(ErrNameNotDeclared, n)
}

func (ns *namespace) resolve(n data.Local) (*Entry, bool) {
	ns.RLock()
	defer ns.RUnlock()
	e, ok := ns.entries[n]
	return e, ok
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
		res.entries[k] = v.snapshot()
	}
	return res
}

func (ns *namespace) Import(imports map[data.Local]*Entry) error {
	ns.Lock()
	defer ns.Unlock()
	names := basics.MapKeys(imports)
	if err := ns.checkDuplicates(names); err != nil {
		return err
	}
	for as, e := range imports {
		res := &Entry{
			name:    as,
			private: true,
			binding: e.binding,
		}
		ns.entries[as] = res
	}
	return nil
}

func (ns *namespace) checkDuplicates(names data.Locals) error {
	duped := data.Locals{}
	for _, n := range names {
		if _, ok := ns.entries[n]; ok {
			duped = append(duped, n)
		}
	}
	if len(duped) > 0 {
		return fmt.Errorf(ErrNameAlreadyDeclared, duped)
	}
	return nil
}

func resolvePublic(
	from Namespace, in Namespace, n data.Local,
) (*Entry, Namespace, error) {
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
