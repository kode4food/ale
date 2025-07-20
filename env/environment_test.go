package env_test

import (
	"fmt"
	"slices"
	"testing"

	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/env"
)

func mustResolveSymbol(ns env.Namespace, s data.Symbol) *env.Entry {
	entry, _, err := env.ResolveSymbol(ns, s)
	if err != nil {
		panic(err)
	}
	return entry
}

func TestResolveSymbol(t *testing.T) {
	as := assert.New(t)

	e := env.NewEnvironment()
	root := e.GetRoot()
	as.NoError(env.BindPublic(root, "public-parent", data.True))
	as.NoError(env.BindPrivate(root, "private-parent", data.True))

	ns := e.GetAnonymous()
	as.NoError(env.BindPublic(ns, "public-child", data.True))
	as.NoError(env.BindPrivate(ns, "private-child", data.True))

	_, _, err := env.ResolveSymbol(ns, LS("public-child"))
	as.NoError(err)

	ent := mustResolveSymbol(ns, LS("private-child"))
	as.NotNil(ent)

	_, _, err = env.ResolveSymbol(ns, LS("public-parent"))
	as.NoError(err)

	ls := LS("private-parent")
	defer as.ExpectPanic(fmt.Errorf(env.ErrNameNotDeclared, ls))
	mustResolveSymbol(ns, ls)
}

func TestResolveValue(t *testing.T) {
	as := assert.New(t)

	e := env.NewEnvironment()
	root := e.GetRoot()
	as.NoError(env.BindPublic(root, "public-parent", data.True))
	as.NoError(env.BindPrivate(root, "private-parent", data.True))

	ns := e.GetAnonymous()
	as.NoError(env.BindPublic(ns, "public-child", data.True))
	as.NoError(env.BindPrivate(ns, "private-child", data.True))

	res, err := env.ResolveValue(ns, LS("public-child"))
	if as.NoError(err) {
		as.True(res)
	}

	as.True(env.MustResolveValue(ns, LS("private-child")))
	res, err = env.ResolveValue(ns, LS("public-parent"))
	if as.NoError(err) {
		as.True(res)
	}

	ls := LS("private-parent")
	defer as.ExpectPanic(fmt.Errorf(env.ErrNameNotDeclared, ls))
	env.MustResolveValue(ns, ls)
}

func TestDomains(t *testing.T) {
	as := assert.New(t)

	e := env.NewEnvironment()
	ns1 := env.MustGetQualified(e, "domain1")
	ns2 := env.MustGetQualified(e, "domain2")

	l := data.Locals{
		e.GetRoot().Domain(),
		ns1.Domain(),
		ns2.Domain(),
	}
	slices.Sort(l)
	r := e.Domains()
	slices.Sort(r)
	as.Equal(l, r)
}
