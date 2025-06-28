package env

import (
	"fmt"
	"sync"

	"github.com/kode4food/ale/internal/basics"
	lang "github.com/kode4food/ale/internal/lang/env"
	"github.com/kode4food/ale/pkg/data"
)

// Environment maintains a mapping of domain names to namespaces
type Environment struct {
	root Namespace
	data map[data.Local]Namespace
	sync.RWMutex
}

const (
	ErrNamespaceNotFound = "namespace not found: %s"
	ErrNamespaceExists   = "namespace already exists: %s"
)

// RootSymbol returns a symbol qualified by the root domain
func RootSymbol(name data.Local) data.Symbol {
	return data.NewQualifiedSymbol(name, lang.RootDomain)
}

// NewEnvironment creates a new synchronous namespace map
func NewEnvironment() *Environment {
	res := &Environment{
		data: map[data.Local]Namespace{},
	}
	res.root = res.newNamespace(lang.RootDomain)
	return res
}

func (e *Environment) Domains() data.Locals {
	e.RLock()
	defer e.RUnlock()
	return append(basics.MapKeys(e.data), lang.RootDomain)
}

func (e *Environment) Snapshot() *Environment {
	e.RLock()
	defer e.RUnlock()
	res := &Environment{
		data: make(map[data.Local]Namespace, len(e.data)),
	}
	res.root = e.root.Snapshot(res)
	for k, v := range e.data {
		res.data[k] = v.Snapshot(res)
	}
	return res
}

// GetRoot returns the root namespace, where built-ins go
func (e *Environment) GetRoot() Namespace {
	return e.root
}

// GetAnonymous returns an anonymous (non-resolvable) namespace
func (e *Environment) GetAnonymous() Namespace {
	return chain(e.root, e.newNamespace(lang.AnonymousDomain))
}

// NewQualified creates a new namespace for the specified domain if it doesn't
// already exist.
func (e *Environment) NewQualified(n data.Local) (Namespace, error) {
	if n == lang.RootDomain {
		return nil, fmt.Errorf(ErrNamespaceExists, lang.RootDomain)
	}
	e.Lock()
	defer e.Unlock()
	if _, ok := e.data[n]; ok {
		return nil, fmt.Errorf(ErrNamespaceExists, n)
	}
	ns := chain(e.root, e.newNamespace(n))
	e.data[n] = ns
	return ns, nil
}

// GetQualified returns the namespace for the specified domain.
func (e *Environment) GetQualified(n data.Local) (Namespace, error) {
	if n == lang.RootDomain {
		return e.root, nil
	}
	e.RLock()
	defer e.RUnlock()
	if ns, ok := e.data[n]; ok {
		return ns, nil
	}
	return nil, fmt.Errorf(ErrNamespaceNotFound, n)
}

func (e *Environment) newNamespace(n data.Local) Namespace {
	return &namespace{
		environment: e,
		entries:     Entries{},
		domain:      n,
	}
}

// ResolveSymbol attempts to resolve a symbol. If it's a qualified symbol, it
// will be retrieved directly from the identified namespace. Otherwise, it will
// be searched in the current namespace
func ResolveSymbol(ns Namespace, s data.Symbol) (*Entry, Namespace, error) {
	if q, ok := s.(data.Qualified); ok {
		e := ns.Environment()
		qns, err := e.GetQualified(q.Domain())
		if err != nil {
			return nil, nil, err
		}
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

// MustGetQualified attempts to retrieve a namespace for the specified domain.
// If the namespace does not exist, it will create it, or panic.
func MustGetQualified(e *Environment, n data.Local) Namespace {
	ns, err := e.GetQualified(n)
	if err == nil {
		return ns
	}
	ns, err = e.NewQualified(n)
	if err != nil {
		panic(err)
	}
	return ns
}
