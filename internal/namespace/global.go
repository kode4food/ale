package namespace

import (
	"fmt"
	"sync"

	"gitlab.com/kode4food/ale/api"
)

type (
	// Manager maintains a mapping of domain names to namespaces
	Manager struct {
		sync.RWMutex
		data map[api.Name]Type
	}

	// Resolver resolves a namespace instance
	Resolver func() Type
)

// Error messages
const (
	SymbolNotBound = "symbol not bound in namespace: %s"
)

const (
	// RootDomain stores built-ins
	RootDomain = api.Name("ale")

	// LocalsDomain is where all local names are interned
	LocalsDomain = api.Name("")

	// UserDomain stores user defined vars
	UserDomain = api.Name("user")

	// AnonymousDomain identifies an anonymous namespace
	AnonymousDomain = api.Name("*anon*")
)

// NewManager creates a new synchronous namespace map
func NewManager() *Manager {
	return &Manager{
		data: map[api.Name]Type{},
	}
}

// Get returns a mapped namespace or instantiates a new one to be cached
func (m *Manager) Get(domain api.Name, res Resolver) Type {
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

// GetRootNamespace returns the root namespace, where built-ins go
func (m *Manager) GetRootNamespace() Type {
	return m.Get(RootDomain, func() Type {
		return m.New(RootDomain)
	})
}

// GetUserNamespace returns the namespace for the user domain
func (m *Manager) GetUserNamespace() Type {
	return m.GetQualified(UserDomain)
}

// GetAnonymousNamespace returns an anonymous (non-resolvable) namespace
func (m *Manager) GetAnonymousNamespace() Type {
	root := m.GetRootNamespace()
	return newChild(root, AnonymousDomain)
}

// GetQualified returns the namespace for the specified domain.
func (m *Manager) GetQualified(n api.Name) Type {
	root := m.GetRootNamespace()
	if n == RootDomain {
		return root
	}
	return m.Get(n, func() Type {
		return newChild(root, n)
	}).(Type)
}

// ResolveSymbol attempts to resolve a symbol. If it's a qualified symbol,
// it will be retrieved directly from the identified namespace. Otherwise
// it will be searched in the current namespace
func ResolveSymbol(ns Type, s api.Symbol) (api.Value, bool) {
	manager := ns.Manager()
	if q, ok := s.(api.QualifiedSymbol); ok {
		qns := manager.GetQualified(q.Domain())
		return qns.Resolve(q.Name())
	}
	return ns.Resolve(s.Name())
}

// MustResolveSymbol attempts to resolve a symbol or explodes violently
func MustResolveSymbol(ns Type, s api.Symbol) api.Value {
	if v, ok := ResolveSymbol(ns, s); ok {
		return v
	}
	panic(fmt.Errorf(SymbolNotBound, s))
}
