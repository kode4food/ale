package namespace_test

import (
	"testing"

	"gitlab.com/kode4food/ale/data"
	"gitlab.com/kode4food/ale/internal/assert"
	"gitlab.com/kode4food/ale/namespace"
)

func TestChaining(t *testing.T) {
	as := assert.New(t)

	manager := namespace.NewManager()
	root := manager.GetRoot()
	root.Bind(data.Name("in-parent"), data.True)

	ns := manager.GetAnonymous()
	ns.Bind(data.Name("in-child"), data.True)

	v1, ok := ns.Resolve(data.Name("in-parent"))
	as.True(ok)
	as.True(v1)

	v2, ok := ns.Resolve(data.Name("in-child"))
	as.True(ok)
	as.True(v2)

	v3, ok := root.Resolve(data.Name("in-child"))
	as.False(ok)
	as.Nil(v3)

	s1 := data.NewLocalSymbol("in-parent")
	v4, ok := namespace.ResolveSymbol(ns, s1)
	as.True(ok)
	as.True(v4)

	v5, ok := namespace.ResolveSymbol(root, s1)
	as.True(ok)
	as.True(v5)

	s2 := data.NewLocalSymbol("in-child")
	v6, ok := namespace.ResolveSymbol(ns, s2)
	as.True(ok)
	as.True(v6)

	v7, ok := namespace.ResolveSymbol(root, s2)
	as.False(ok)
	as.Nil(v7)

	s3 := namespace.RootSymbol("in-parent")
	v8, ok := namespace.ResolveSymbol(ns, s3)
	as.True(ok)
	as.True(v8)
}
