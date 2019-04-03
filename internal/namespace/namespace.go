package namespace

import (
	"fmt"
	"sync"

	"gitlab.com/kode4food/ale/api"
)

type (
	// Type represents a namespace
	Type interface {
		Manager() *Manager
		Domain() api.Name
		Resolve(api.Name) (api.Value, bool)
		IsDeclared(api.Name) bool
		Declare(api.Name)
		IsBound(api.Name) bool
		Bind(api.Name, api.Value)
	}

	namespace struct {
		manager  *Manager
		domain   api.Name
		entries  entries
		entMutex sync.RWMutex
	}

	entry struct {
		bound bool
		value api.Value
	}

	entries map[api.Name]*entry
)

// Error messages
const (
	NameAlreadyBound = "name is already bound in namespace: %s"
)

// New constructs a new namespace
func (m *Manager) New(n api.Name) Type {
	return &namespace{
		manager: m,
		entries: entries{},
		domain:  n,
	}
}

func (ns *namespace) Manager() *Manager {
	return ns.manager
}

func (ns *namespace) Resolve(n api.Name) (api.Value, bool) {
	ns.entMutex.RLock()
	defer ns.entMutex.RUnlock()
	if res, ok := ns.entries[n]; ok {
		return res.value, res.bound
	}
	return api.Nil, false
}

func (ns *namespace) IsDeclared(n api.Name) bool {
	ns.entMutex.RLock()
	defer ns.entMutex.RUnlock()
	_, ok := ns.entries[n]
	return ok
}

func (ns *namespace) Declare(n api.Name) {
	ns.entMutex.Lock()
	defer ns.entMutex.Unlock()
	if _, ok := ns.entries[n]; !ok {
		ns.entries[n] = &entry{
			bound: false,
			value: api.Nil,
		}
	}
}

func (ns *namespace) IsBound(n api.Name) bool {
	ns.entMutex.RLock()
	defer ns.entMutex.RUnlock()
	e, ok := ns.entries[n]
	return ok && e.bound
}

func (ns *namespace) Bind(n api.Name, v api.Value) {
	ns.entMutex.Lock()
	defer ns.entMutex.Unlock()
	e, ok := ns.entries[n]
	if !ok || !e.bound {
		ns.entries[n] = &entry{
			bound: true,
			value: v,
		}
		return
	}
	panic(api.String(fmt.Sprintf(NameAlreadyBound, n)))
}

func (ns *namespace) Domain() api.Name {
	return ns.domain
}
