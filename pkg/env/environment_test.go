package env_test

import (
	"fmt"
	"testing"

	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/pkg/data"
	"github.com/kode4food/ale/pkg/env"
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
	as.Nil(env.BindPublic(root, "public-parent", data.True))
	as.Nil(env.BindPrivate(root, "private-parent", data.True))

	ns := e.GetAnonymous()
	as.Nil(env.BindPublic(ns, "public-child", data.True))
	as.Nil(env.BindPrivate(ns, "private-child", data.True))

	_, _, err := env.ResolveSymbol(ns, LS("public-child"))
	as.Nil(err)

	ent := mustResolveSymbol(ns, LS("private-child"))
	as.NotNil(ent)

	_, _, err = env.ResolveSymbol(ns, LS("public-parent"))
	as.Nil(err)

	ls := LS("private-parent")
	defer as.ExpectPanic(fmt.Errorf(env.ErrNameNotDeclared, ls))
	mustResolveSymbol(ns, ls)
}

func TestResolveValue(t *testing.T) {
	as := assert.New(t)

	e := env.NewEnvironment()
	root := e.GetRoot()
	as.Nil(env.BindPublic(root, "public-parent", data.True))
	as.Nil(env.BindPrivate(root, "private-parent", data.True))

	ns := e.GetAnonymous()
	as.Nil(env.BindPublic(ns, "public-child", data.True))
	as.Nil(env.BindPrivate(ns, "private-child", data.True))

	res, err := env.ResolveValue(ns, LS("public-child"))
	as.True(res)
	as.Nil(err)

	as.True(env.MustResolveValue(ns, LS("private-child")))
	res, err = env.ResolveValue(ns, LS("public-parent"))
	as.True(res)
	as.Nil(err)

	ls := LS("private-parent")
	defer as.ExpectPanic(fmt.Errorf(env.ErrNameNotDeclared, ls))
	env.MustResolveValue(ns, ls)
}

func TestDomains(t *testing.T) {
	as := assert.New(t)

	e := env.NewEnvironment()
	ns1 := e.GetQualified("domain1")
	ns2 := e.GetQualified("domain2")

	as.Equal(data.Locals{
		e.GetRoot().Domain(),
		ns1.Domain(),
		ns2.Domain(),
	}, e.Domains().Sorted())
}
