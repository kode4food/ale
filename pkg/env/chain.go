package env

import "github.com/kode4food/ale/pkg/data"

type chainedNamespace struct {
	child  Namespace
	parent Namespace
}

func newChild(parent Namespace, n data.Local) Namespace {
	e := parent.Environment()
	return chain(parent, e.New(n))
}

func chain(parent Namespace, child Namespace) *chainedNamespace {
	return &chainedNamespace{
		parent: parent,
		child:  child,
	}
}

func (ns *chainedNamespace) Snapshot(e *Environment) (Namespace, error) {
	p, err := ns.parent.Snapshot(e)
	if err != nil {
		return nil, err
	}
	c, err := ns.child.Snapshot(e)
	if err != nil {
		return nil, err
	}
	return &chainedNamespace{
		parent: p,
		child:  c,
	}, nil
}

func (ns *chainedNamespace) Environment() *Environment {
	return ns.child.Environment()
}

func (ns *chainedNamespace) Domain() data.Local {
	return ns.child.Domain()
}

func (ns *chainedNamespace) Declared() data.Locals {
	return ns.child.Declared()
}

func (ns *chainedNamespace) Public(n data.Local) (Entry, error) {
	return ns.child.Public(n)
}

func (ns *chainedNamespace) Private(n data.Local) (Entry, error) {
	return ns.child.Private(n)
}

func (ns *chainedNamespace) Resolve(n data.Local) (Entry, Namespace, error) {
	if e, in, err := ns.child.Resolve(n); err == nil {
		return e, in, nil
	}
	return resolvePublic(ns, ns.parent, n)
}
