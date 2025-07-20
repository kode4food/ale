package env_test

import (
	"fmt"
	"slices"
	"testing"

	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	lang "github.com/kode4food/ale/internal/lang/env"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/env"
)

func TestDeclarations(t *testing.T) {
	as := assert.New(t)

	e := env.NewEnvironment()
	root := e.GetRoot()
	as.Equal(e, root.Environment())
	as.Equal(lang.RootDomain, root.Domain())

	as.NoError(env.BindPublic(root, "public2", data.True))
	as.NoError(env.BindPublic(root, "public1", data.True))
	as.NoError(env.BindPrivate(root, "private", data.True))

	n := root.Declared()
	slices.Sort(n)
	as.Equal(2, len(n))
	as.Equal(LS("public1"), n[0])
	as.Equal(LS("public2"), n[1])

	e2, in, err := root.Resolve(n[0])
	if as.NoError(err) {
		as.NotNil(e2)
		as.NotNil(in)
	}

	as.Equal(n[0], e2.Name())
	as.Equal(root, in)

	e3, err := root.Public(n[0])
	if as.NoError(err) {
		as.Equal(e2, e3)
	}
}

func TestChaining(t *testing.T) {
	as := assert.New(t)

	e := env.NewEnvironment()
	root := e.GetRoot()
	as.NoError(env.BindPublic(root, "in-parent", data.True))

	ns := e.GetAnonymous()
	as.NoError(env.BindPublic(ns, "in-child", data.True))

	as.True(as.IsBound(ns, "in-parent"))
	as.True(as.IsBound(ns, "in-child"))
	as.True(as.IsBound(root, "in-parent"))
	as.IsNotDeclared(root, "in-child")
	s3 := env.RootSymbol("in-parent")
	v8, err := env.ResolveValue(ns, s3)
	if as.NoError(err) {
		as.True(v8)
	}
}

func TestBinding(t *testing.T) {
	as := assert.New(t)

	e := env.NewEnvironment()
	root := e.GetRoot()
	d, err := root.Public("some-name")
	if as.NoError(err) {
		v, err := d.Value()
		as.Nil(v)
		as.EqualError(err, fmt.Sprintf(env.ErrNameNotBound, d.Name()))

		err = d.Bind(S("some-value"))
		if as.NoError(err) {
			err = d.Bind(S("some-other-value"))
			as.EqualError(err, fmt.Sprintf(env.ErrNameAlreadyBound, d.Name()))
		}

		v, err = d.Value()
		if as.NoError(err) {
			as.String("some-value", v)
		}
	}
}

func TestRedeclaration(t *testing.T) {
	as := assert.New(t)
	e := env.NewEnvironment()
	root := e.GetRoot()
	d1, err := root.Public("some-name")
	if as.NoError(err) {
		as.NotNil(d1)
	}

	d2, err := root.Public("some-name")
	if as.NoError(err) {
		as.NotNil(d2)
		as.Equal(d1, d2)
	}

	_, err = root.Private("some-name")
	as.EqualError(err, fmt.Sprintf(env.ErrNameAlreadyDeclared, "some-name"))

	d3, err := root.Private("other-name")
	if as.NoError(err) {
		as.NotNil(d3)
	}

	d4, err := root.Private("other-name")
	if as.NoError(err) {
		as.NotNil(d4)
		as.Equal(d3, d4)
	}

	_, err = root.Public("other-name")
	as.EqualError(err, fmt.Sprintf(env.ErrNameAlreadyDeclared, "other-name"))
}
