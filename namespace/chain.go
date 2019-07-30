package namespace

import "github.com/kode4food/ale/data"

type chainedNamespace struct {
	child  Type
	parent Type
}

func newChild(parent Type, n data.Name) Type {
	m := parent.Manager()
	return chain(parent, m.New(n))
}

func chain(parent Type, child Type) *chainedNamespace {
	return &chainedNamespace{
		parent: parent,
		child:  child,
	}
}

func (ns *chainedNamespace) Manager() *Manager {
	return ns.child.Manager()
}

func (ns *chainedNamespace) Domain() data.Name {
	return ns.child.Domain()
}

func (ns *chainedNamespace) Declare(n data.Name) Entry {
	return ns.child.Declare(n)
}

func (ns *chainedNamespace) Resolve(n data.Name) (Entry, bool) {
	if e, ok := ns.child.Resolve(n); ok {
		return e, true
	}
	return ns.parent.Resolve(n)
}
