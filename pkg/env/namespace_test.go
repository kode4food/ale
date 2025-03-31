package env_test

import (
	"fmt"
	"testing"

	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/pkg/data"
	"github.com/kode4food/ale/pkg/env"
)

func TestDeclarations(t *testing.T) {
	as := assert.New(t)

	e := env.NewEnvironment()
	root := e.GetRoot()
	as.Equal(e, root.Environment())
	as.Equal(env.RootDomain, root.Domain())

	as.Nil(env.BindPublic(root, "public2", data.True))
	as.Nil(env.BindPublic(root, "public1", data.True))
	as.Nil(env.BindPrivate(root, "private", data.True))

	n := root.Declared()
	as.Equal(2, len(n))
	as.Equal(LS("public1"), n[0])
	as.Equal(LS("public2"), n[1])

	e2, in, err := root.Resolve(n[0])
	as.NotNil(e2)
	as.NotNil(in)
	as.Nil(err)

	as.Equal(n[0], e2.Name())
	as.Equal(root, in)

	e3, err := root.Public(n[0])
	as.Equal(e2, e3)
	as.Nil(err)
}

func TestChaining(t *testing.T) {
	as := assert.New(t)

	e := env.NewEnvironment()
	root := e.GetRoot()
	as.Nil(env.BindPublic(root, "in-parent", data.True))

	ns := e.GetAnonymous()
	as.Nil(env.BindPublic(ns, "in-child", data.True))

	as.True(as.IsBound(ns, "in-parent"))
	as.True(as.IsBound(ns, "in-child"))
	as.True(as.IsBound(root, "in-parent"))
	as.IsNotDeclared(root, "in-child")
	s3 := env.RootSymbol("in-parent")
	v8, err := env.ResolveValue(ns, s3)
	as.True(v8)
	as.Nil(err)
}

func TestBinding(t *testing.T) {
	as := assert.New(t)

	e := env.NewEnvironment()
	root := e.GetRoot()
	d, err := root.Public("some-name")
	as.Nil(err)

	v, err := d.Value()
	as.Nil(v)
	as.EqualError(err, fmt.Sprintf(env.ErrNameNotBound, d.Name()))

	err = d.Bind(S("some-value"))
	as.Nil(err)
	err = d.Bind(S("some-other-value"))
	as.EqualError(err, fmt.Sprintf(env.ErrNameAlreadyBound, d.Name()))

	v, err = d.Value()
	as.Nil(err)
	as.String("some-value", v)
}
