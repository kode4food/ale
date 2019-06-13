package namespace

import (
	"fmt"
	"sync"

	"gitlab.com/kode4food/ale/data"
)

type (
	// Manager maintains a mapping of domain names to namespaces
	Manager struct {
		sync.RWMutex
		data map[data.Name]Type
	}

	// Resolver resolves a namespace instance
	Resolver func() Type
)

// Error messages
const (
	SymbolNotDeclared = "symbol not declared in namespace: %s"
	SymbolNotBound    = "symbol not bound in namespace: %s"
)

const (
	// RootDomain stores built-ins
	RootDomain = data.Name("ale")

	// AnonymousDomain identifies an anonymous namespace
	AnonymousDomain = data.Name("*anon*")
)

// RootSymbol returns a symbol qualified by the root domain
func RootSymbol(name data.Name) data.Symbol {
	return data.NewQualifiedSymbol(name, RootDomain)
}

// NewManager creates a new synchronous namespace map
func NewManager() *Manager {
	return &Manager{
		data: map[data.Name]Type{},
	}
}

// New constructs a new namespace
func (m *Manager) New(n data.Name) Type {
	return &namespace{
		manager: m,
		entries: entries{},
		domain:  n,
	}
}

// Get returns a mapped namespace or instantiates a new one to be cached
func (m *Manager) Get(domain data.Name, res Resolver) Type {
	m.RLock()
	r, ok := m.data[domain]
	m.RUnlock()
	if ok {
		return r
	}

	r = res()
	m.Lock()
	defer m.Unlock()
	if orig, ok := m.data[domain]; ok {
		return orig
	}
	m.data[domain] = r
	return r
}

// GetRoot returns the root namespace, where built-ins go
func (m *Manager) GetRoot() Type {
	return m.Get(RootDomain, func() Type {
		return m.New(RootDomain)
	})
}

// GetAnonymous returns an anonymous (non-resolvable) namespace
func (m *Manager) GetAnonymous() Type {
	root := m.GetRoot()
	return chain(root, &anonymous{
		Type: m.New(AnonymousDomain),
	})
}

// GetQualified returns the namespace for the specified domain.
func (m *Manager) GetQualified(n data.Name) Type {
	root := m.GetRoot()
	if n == RootDomain {
		return root
	}
	return m.Get(n, func() Type {
		return newChild(root, n)
	})
}

// ResolveSymbol attempts to resolve a symbol. If it's a qualified symbol,
// it will be retrieved directly from the identified namespace. Otherwise
// it will be searched in the current namespace
func ResolveSymbol(ns Type, s data.Symbol) (Entry, bool) {
	if q, ok := s.(data.QualifiedSymbol); ok {
		manager := ns.Manager()
		qns := manager.GetQualified(q.Domain())
		return qns.Resolve(q.Name())
	}
	return ns.Resolve(s.Name())
}

// MustResolveSymbol attempts to resolve a symbol or explodes violently
func MustResolveSymbol(ns Type, s data.Symbol) Entry {
	if entry, ok := ResolveSymbol(ns, s); ok {
		return entry
	}
	panic(fmt.Errorf(SymbolNotDeclared, s.Name()))
}

// ResolveValue attempts to resolve a symbol to a bound value
func ResolveValue(ns Type, s data.Symbol) (data.Value, bool) {
	if e, ok := ResolveSymbol(ns, s); ok && e.IsBound() {
		return e.Value(), true
	}
	return data.Null, false
}

// MustResolveValue attempts to resolve a value or explodes violently
func MustResolveValue(ns Type, s data.Symbol) data.Value {
	if v, ok := ResolveValue(ns, s); ok {
		return v
	}
	panic(fmt.Errorf(SymbolNotBound, s))
}
