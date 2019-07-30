package namespace_test

import (
	"testing"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/assert"
	"github.com/kode4food/ale/namespace"
)

func TestChaining(t *testing.T) {
	as := assert.New(t)

	manager := namespace.NewManager()
	root := manager.GetRoot()
	root.Declare(data.Name("in-parent")).Bind(data.True)

	ns := manager.GetAnonymous()
	ns.Declare(data.Name("in-child")).Bind(data.True)

	e1, ok := ns.Resolve(data.Name("in-parent"))
	as.True(ok && e1.IsBound())
	as.True(e1.Value())

	e2, ok := ns.Resolve(data.Name("in-child"))
	as.True(ok && e2.IsBound())
	as.True(e2.Value())

	e3, ok := root.Resolve(data.Name("in-child"))
	as.False(ok)
	as.Nil(e3)

	s1 := data.NewLocalSymbol("in-parent")
	v4, ok := namespace.ResolveValue(ns, s1)
	as.True(ok)
	as.True(v4)

	v5, ok := namespace.ResolveValue(root, s1)
	as.True(ok)
	as.True(v5)

	s2 := data.NewLocalSymbol("in-child")
	v6, ok := namespace.ResolveValue(ns, s2)
	as.True(ok)
	as.True(v6)

	v7, ok := namespace.ResolveValue(root, s2)
	as.False(ok)
	as.Nil(v7)

	s3 := namespace.RootSymbol("in-parent")
	v8, ok := namespace.ResolveValue(ns, s3)
	as.True(ok)
	as.True(v8)
}
