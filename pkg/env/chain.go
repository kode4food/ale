package env

import "github.com/kode4food/ale/pkg/data"

type chainedNamespace struct {
	Namespace // child
	parent    Namespace
}

func chain(parent Namespace, child Namespace) *chainedNamespace {
	return &chainedNamespace{
		Namespace: child,
		parent:    parent,
	}
}

func (ns *chainedNamespace) Snapshot(e *Environment) Namespace {
	return &chainedNamespace{
		Namespace: ns.Namespace.Snapshot(e),
		parent:    ns.parent.Snapshot(e),
	}
}

func (ns *chainedNamespace) Resolve(n data.Local) (*Entry, Namespace, error) {
	if e, in, err := ns.Namespace.Resolve(n); err == nil {
		return e, in, nil
	}
	return resolvePublic(ns, ns.parent, n)
}
