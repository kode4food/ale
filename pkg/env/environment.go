package env

import (
	"sync"

	"github.com/kode4food/ale/internal/basics"
	lang "github.com/kode4food/ale/internal/lang/env"
	"github.com/kode4food/ale/pkg/data"
)

type (
	// Environment maintains a mapping of domain names to namespaces
	Environment struct {
		data map[data.Local]Namespace
		sync.RWMutex
	}

	// Resolver resolves a namespace instance
	Resolver func() Namespace
)

// RootSymbol returns a symbol qualified by the root domain
func RootSymbol(name data.Local) data.Symbol {
	return data.NewQualifiedSymbol(name, lang.RootDomain)
}

// NewEnvironment creates a new synchronous namespace map
func NewEnvironment() *Environment {
	return &Environment{
		data: map[data.Local]Namespace{},
	}
}

func (e *Environment) Domains() data.Locals {
	e.RLock()
	defer e.RUnlock()
	return basics.MapKeys(e.data)
}

func (e *Environment) Snapshot() *Environment {
	e.RLock()
	defer e.RUnlock()
	res := &Environment{
		data: make(map[data.Local]Namespace, len(e.data)),
	}
	for k, v := range e.data {
		res.data[k] = v.Snapshot(res)
	}
	return res
}

// Get returns a mapped namespace or instantiates a new one to be cached
func (e *Environment) Get(domain data.Local, res Resolver) Namespace {
	if r, ok := e.get(domain); ok {
		return r
	}
	e.Lock()
	defer e.Unlock()
	if r, ok := e.data[domain]; ok {
		return r
	}
	r := res()
	e.data[domain] = r
	return r
}

func (e *Environment) get(domain data.Local) (Namespace, bool) {
	e.RLock()
	defer e.RUnlock()
	r, ok := e.data[domain]
	return r, ok
}

// GetRoot returns the root namespace, where built-ins go
func (e *Environment) GetRoot() Namespace {
	return e.Get(lang.RootDomain, func() Namespace {
		return e.newNamespace(lang.RootDomain)
	})
}

// GetAnonymous returns an anonymous (non-resolvable) namespace
func (e *Environment) GetAnonymous() Namespace {
	return chain(e.GetRoot(), e.newNamespace(lang.AnonymousDomain))
}

// GetQualified returns the namespace for the specified domain.
func (e *Environment) GetQualified(n data.Local) Namespace {
	root := e.GetRoot()
	if n == lang.RootDomain {
		return root
	}
	return e.Get(n, func() Namespace {
		return chain(root, e.newNamespace(n))
	})
}

func (e *Environment) newNamespace(n data.Local) Namespace {
	return &namespace{
		environment: e,
		entries:     entries{},
		domain:      n,
	}
}

// ResolveSymbol attempts to resolve a symbol. If it's a qualified symbol, it
// will be retrieved directly from the identified namespace. Otherwise, it will
// be searched in the current namespace
func ResolveSymbol(ns Namespace, s data.Symbol) (*Entry, Namespace, error) {
	if q, ok := s.(data.Qualified); ok {
		e := ns.Environment()
		qns := e.GetQualified(q.Domain())
		return resolvePublic(ns, qns, q.Name())
	}
	return ns.Resolve(s.Name())
}

// ResolveValue attempts to resolve a symbol to a bound value
func ResolveValue(ns Namespace, s data.Symbol) (data.Value, error) {
	e, _, err := ResolveSymbol(ns, s)
	if err != nil {
		return nil, err
	}
	return e.Value()
}

// MustResolveValue attempts to resolve a value or explodes violently
func MustResolveValue(ns Namespace, s data.Symbol) data.Value {
	v, err := ResolveValue(ns, s)
	if err != nil {
		panic(err)
	}
	return v
}
