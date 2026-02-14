package env_test

import (
	"fmt"
	"slices"
	"testing"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/env"
	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	lang "github.com/kode4food/ale/internal/lang/env"
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
	if as.NoError(err) && as.NotNil(e2) && as.NotNil(in) {
		as.Equal(n[0], e2.Name())
		as.Equal(root, in)

		e3, err := root.Public(n[0])
		if as.NoError(err) {
			as.Equal(e2, e3)
		}
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
	if as.NoError(err) && as.NotNil(d1) {
		d2, err := root.Public("some-name")
		if as.NoError(err) && as.NotNil(d2) {
			as.Equal(d1, d2)
		}
	}

	_, err = root.Private("some-name")
	as.EqualError(err, fmt.Sprintf(env.ErrNameAlreadyDeclared, "some-name"))

	d3, err := root.Private("other-name")
	if as.NoError(err) && as.NotNil(d3) {
		d4, err := root.Private("other-name")
		if as.NoError(err) && as.NotNil(d4) {
			as.Equal(d3, d4)
		}
	}

	_, err = root.Public("other-name")
	as.EqualError(err, fmt.Sprintf(env.ErrNameAlreadyDeclared, "other-name"))
}

func TestImport(t *testing.T) {
	as := assert.New(t)
	e := env.NewEnvironment()

	src, err := e.NewQualified("src")
	if as.NoError(err) {
		as.NoError(env.BindPublic(src, "public-value", I(42)))
		as.NoError(env.BindPrivate(src, "private-value", I(99)))
	}

	dst, err := e.NewQualified("dst")
	if as.NoError(err) {
		pub, _, err := src.Resolve("public-value")
		if as.NoError(err) {
			as.NoError(dst.Import(env.Entries{
				"alias": pub,
			}))
		}
	}

	d := dst.Declared()
	as.Equal(0, len(d))

	e1, _, err := dst.Resolve("alias")
	if as.NoError(err) {
		v, err := e1.Value()
		if as.NoError(err) {
			as.Equal(I(42), v)
		}
		as.True(e1.IsPrivate())
	}
}

func TestImportDuplicatesAndResolvePublic(t *testing.T) {
	as := assert.New(t)
	e := env.NewEnvironment()

	from, err := e.NewQualified("from")
	if as.NoError(err) {
		as.NoError(env.BindPublic(from, "name", I(1)))
	}

	to, err := e.NewQualified("to")
	if as.NoError(err) {
		as.NoError(env.BindPublic(to, "name", I(2)))
	}

	fEntry, _, err := from.Resolve("name")
	if as.NoError(err) {
		err = to.Import(env.Entries{
			"name": fEntry,
		})
		as.EqualError(err, fmt.Sprintf(env.ErrNameAlreadyDeclared, "[name]"))
	}

	pEntry, _, err := from.Resolve("name")
	if as.NoError(err) {
		as.NoError(to.Import(env.Entries{
			"private-name": pEntry,
		}))
	}

	q := QS("to", "private-name")
	_, _, err = env.ResolveSymbol(e.GetRoot(), q)
	as.EqualError(err, fmt.Sprintf(env.ErrNameNotDeclared, "private-name"))

	ent, _, err := env.ResolveSymbol(to, q)
	if as.NoError(err) && as.NotNil(ent) {
		v, err := ent.Value()
		if as.NoError(err) {
			as.Equal(I(1), v)
		}
	}
}
