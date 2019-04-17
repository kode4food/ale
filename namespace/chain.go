package namespace

import "gitlab.com/kode4food/ale/api"

type chainedNamespace struct {
	child  Type
	parent Type
}

func newChild(parent Type, n api.Name) Type {
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

func (ns *chainedNamespace) Domain() api.Name {
	return ns.child.Domain()
}

func (ns *chainedNamespace) In(n api.Name) (Type, bool) {
	if ns, ok := ns.child.In(n); ok {
		return ns, true
	}
	return ns.parent.In(n)
}

func (ns *chainedNamespace) Resolve(n api.Name) (api.Value, bool) {
	if ns.child.IsDeclared(n) {
		return ns.child.Resolve(n)
	}
	return ns.parent.Resolve(n)
}

func (ns *chainedNamespace) IsDeclared(n api.Name) bool {
	return ns.child.IsDeclared(n)
}

func (ns *chainedNamespace) Declare(n api.Name) {
	ns.child.Declare(n)
}

func (ns *chainedNamespace) IsBound(n api.Name) bool {
	return ns.child.IsBound(n)
}

func (ns *chainedNamespace) Bind(n api.Name, v api.Value) {
	ns.child.Bind(n, v)
}
