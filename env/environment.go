package env

import (
	"fmt"
	"sync"

	"github.com/kode4food/ale/data"
)

type (
	// Environment maintains a mapping of domain names to namespaces
	Environment struct {
		sync.Mutex
		data map[data.Local]Namespace
	}

	// Resolver resolves a namespace instance
	Resolver func() Namespace
)

// Error messages
const (
	ErrSymbolNotDeclared = "symbol not declared in namespace: %s"
	ErrSymbolNotBound    = "symbol not bound in namespace: %s"
)

const (
	// RootDomain stores built-ins
	RootDomain = data.Local("ale")

	// AnonymousDomain identifies an anonymous namespace
	AnonymousDomain = data.Local("*anon*")
)

// RootSymbol returns a symbol qualified by the root domain
func RootSymbol(name data.Local) data.Symbol {
	return data.NewQualifiedSymbol(name, RootDomain)
}

// NewEnvironment creates a new synchronous namespace map
func NewEnvironment() *Environment {
	return &Environment{
		data: map[data.Local]Namespace{},
	}
}

func (e *Environment) Domains() data.Locals {
	res := make(data.Locals, 0, len(e.data))
	for k := range e.data {
		res = append(res, k)
	}
	return res.Sorted()
}

func (e *Environment) Snapshot() (*Environment, error) {
	e.Lock()
	defer e.Unlock()

	res := &Environment{
		data: make(map[data.Local]Namespace, len(e.data)),
	}
	for k, v := range e.data {
		s, err := v.Snapshot(res)
		if err != nil {
			return nil, err
		}
		res.data[k] = s
	}
	return res, nil
}

// New constructs a new namespace
func (e *Environment) New(n data.Local) Namespace {
	return &namespace{
		environment: e,
		entries:     entries{},
		domain:      n,
	}
}

// Get returns a mapped namespace or instantiates a new one to be cached
func (e *Environment) Get(domain data.Local, res Resolver) Namespace {
	e.Lock()
	defer e.Unlock()
	if r, ok := e.data[domain]; ok {
		return r
	}
	r := res()
	e.data[domain] = r
	return r
}

// GetRoot returns the root namespace, where built-ins go
func (e *Environment) GetRoot() Namespace {
	return e.Get(RootDomain, func() Namespace {
		return e.New(RootDomain)
	})
}

// GetAnonymous returns an anonymous (non-resolvable) namespace
func (e *Environment) GetAnonymous() Namespace {
	root := e.GetRoot()
	return chain(root, &anonymous{
		Namespace: e.New(AnonymousDomain),
	})
}

// GetQualified returns the namespace for the specified domain.
func (e *Environment) GetQualified(n data.Local) Namespace {
	root := e.GetRoot()
	if n == RootDomain {
		return root
	}
	return e.Get(n, func() Namespace {
		return newChild(root, n)
	})
}

// ResolveSymbol attempts to resolve a symbol. If it's a qualified symbol, it
// will be retrieved directly from the identified namespace. Otherwise it will
// be searched in the current namespace
func ResolveSymbol(ns Namespace, s data.Symbol) (Entry, bool) {
	if q, ok := s.(data.Qualified); ok {
		e := ns.Environment()
		qns := e.GetQualified(q.Domain())
		return resolvePublic(ns, qns, q.Name())
	}
	return ns.Resolve(s.Name())
}

// MustResolveSymbol attempts to resolve a symbol or explodes violently
func MustResolveSymbol(ns Namespace, s data.Symbol) Entry {
	if entry, ok := ResolveSymbol(ns, s); ok {
		return entry
	}
	panic(fmt.Errorf(ErrSymbolNotDeclared, s.Name()))
}

// ResolveValue attempts to resolve a symbol to a bound value
func ResolveValue(ns Namespace, s data.Symbol) (data.Value, bool) {
	if e, ok := ResolveSymbol(ns, s); ok && e.IsBound() {
		return e.Value(), true
	}
	return data.Nil, false
}

// MustResolveValue attempts to resolve a value or explodes violently
func MustResolveValue(ns Namespace, s data.Symbol) data.Value {
	if v, ok := ResolveValue(ns, s); ok {
		return v
	}
	panic(fmt.Errorf(ErrSymbolNotBound, s))
}
