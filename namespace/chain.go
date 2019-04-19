package namespace

import "gitlab.com/kode4food/ale/data"

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

func (ns *chainedNamespace) In(n data.Name) (Type, bool) {
	if ns, ok := ns.child.In(n); ok {
		return ns, true
	}
	return ns.parent.In(n)
}

func (ns *chainedNamespace) Resolve(n data.Name) (data.Value, bool) {
	if ns.child.IsDeclared(n) {
		return ns.child.Resolve(n)
	}
	return ns.parent.Resolve(n)
}

func (ns *chainedNamespace) IsDeclared(n data.Name) bool {
	return ns.child.IsDeclared(n)
}

func (ns *chainedNamespace) Declare(n data.Name) {
	ns.child.Declare(n)
}

func (ns *chainedNamespace) IsBound(n data.Name) bool {
	return ns.child.IsBound(n)
}

func (ns *chainedNamespace) Bind(n data.Name, v data.Value) {
	ns.child.Bind(n, v)
}
