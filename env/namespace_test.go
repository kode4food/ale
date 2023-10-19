package env_test

import (
	"fmt"
	"testing"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/env"
	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
)

func TestDeclarations(t *testing.T) {
	as := assert.New(t)

	e := env.NewEnvironment()
	root := e.GetRoot()
	as.Equal(e, root.Environment())
	as.Equal(env.RootDomain, root.Domain())

	root.Declare("public2").Bind(data.True)
	root.Private("private").Bind(data.True)
	root.Declare("public1").Bind(data.True)

	n := root.Declared()
	as.Equal(2, len(n))
	as.Equal(LS("public1"), n[0])
	as.Equal(LS("public2"), n[1])

	e1, ok := root.Resolve(n[0])
	as.NotNil(e1)
	as.True(ok)

	as.Equal(n[0], e1.Name())
	as.Equal(root, e1.Owner())

	e2 := root.Declare(n[0])
	as.Equal(e1, e2)
}

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

func TestBinding(t *testing.T) {
	as := assert.New(t)

	e := env.NewEnvironment()
	root := e.GetRoot()
	d := root.Declare("some-name")

	func() {
		defer as.ExpectPanic(fmt.Errorf(env.ErrNameNotBound, d.Name()))
		d.Value()
	}()

	d.Bind(S("some-value"))

	func() {
		defer as.ExpectPanic(fmt.Errorf(env.ErrNameAlreadyBound, d.Name()))
		d.Bind(S("some-other-value"))
	}()

	as.String("some-value", d.Value())
}
