package env_test

import (
	"testing"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/env"
	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
)

func TestChaining(t *testing.T) {
	as := assert.New(t)

	e := env.NewEnvironment()
	root := e.GetRoot()
	root.Declare("in-parent").Bind(data.True)

	ns := e.GetAnonymous()
	ns.Declare("in-child").Bind(data.True)

	e1, ok := ns.Resolve("in-parent")
	as.True(ok && e1.IsBound())
	as.True(e1.Value())

	e2, ok := ns.Resolve("in-child")
	as.True(ok && e2.IsBound())
	as.True(e2.Value())

	e3, ok := root.Resolve("in-child")
	as.False(ok)
	as.Nil(e3)

	s1 := LS("in-parent")
	v4, ok := env.ResolveValue(ns, s1)
	as.True(ok)
	as.True(v4)

	v5, ok := env.ResolveValue(root, s1)
	as.True(ok)
	as.True(v5)

	s2 := LS("in-child")
	v6, ok := env.ResolveValue(ns, s2)
	as.True(ok)
	as.True(v6)

	v7, ok := env.ResolveValue(root, s2)
	as.False(ok)
	as.Nil(v7)

	s3 := env.RootSymbol("in-parent")
	v8, ok := env.ResolveValue(ns, s3)
	as.True(ok)
	as.True(v8)
}
