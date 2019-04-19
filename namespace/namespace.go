package namespace

import (
	"fmt"
	"sync"

	"gitlab.com/kode4food/ale/data"
)

type (
	// Type represents a namespace
	Type interface {
		Manager() *Manager
		Domain() data.Name
		Resolve(data.Name) (data.Value, bool)
		In(data.Name) (Type, bool)
		IsDeclared(data.Name) bool
		Declare(data.Name)
		IsBound(data.Name) bool
		Bind(data.Name, data.Value)
	}

	namespace struct {
		manager  *Manager
		domain   data.Name
		entries  entries
		entMutex sync.RWMutex
	}

	anonymous struct {
		Type
	}

	entry struct {
		bound bool
		value data.Value
	}

	entries map[data.Name]*entry
)

// Error messages
const (
	NameAlreadyBound = "name is already bound in namespace: %s"
)

func (ns *namespace) Manager() *Manager {
	return ns.manager
}

func (ns *namespace) Resolve(n data.Name) (data.Value, bool) {
	ns.entMutex.RLock()
	defer ns.entMutex.RUnlock()
	if res, ok := ns.entries[n]; ok {
		return res.value, res.bound
	}
	return data.Nil, false
}

func (ns *namespace) In(n data.Name) (Type, bool) {
	ns.entMutex.RLock()
	defer ns.entMutex.RUnlock()
	if _, ok := ns.entries[n]; ok {
		return ns, true
	}
	return nil, false
}

func (ns *namespace) IsDeclared(n data.Name) bool {
	ns.entMutex.RLock()
	defer ns.entMutex.RUnlock()
	_, ok := ns.entries[n]
	return ok
}

func (ns *namespace) Declare(n data.Name) {
	ns.entMutex.Lock()
	defer ns.entMutex.Unlock()
	if _, ok := ns.entries[n]; !ok {
		ns.entries[n] = &entry{
			bound: false,
			value: data.Nil,
		}
	}
}

func (ns *namespace) IsBound(n data.Name) bool {
	ns.entMutex.RLock()
	defer ns.entMutex.RUnlock()
	e, ok := ns.entries[n]
	return ok && e.bound
}

func (ns *namespace) Bind(n data.Name, v data.Value) {
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
	panic(data.String(fmt.Sprintf(NameAlreadyBound, n)))
}

func (ns *namespace) Domain() data.Name {
	return ns.domain
}

func (*anonymous) In(data.Name) (Type, bool) {
	return nil, false
}
