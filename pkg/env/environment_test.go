package env_test

import (
	"fmt"
	"testing"

	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/pkg/data"
	"github.com/kode4food/ale/pkg/env"
)

func mustResolveSymbol(ns env.Namespace, s data.Symbol) env.Entry {
	entry, err := env.ResolveSymbol(ns, s)
	if err != nil {
		panic(err)
	}
	return entry
}

func TestResolveSymbol(t *testing.T) {
	as := assert.New(t)

	e := env.NewEnvironment()
	root := e.GetRoot()
	as.Nil(root.Declare("public-parent").Bind(data.True))
	as.Nil(root.Private("private-parent").Bind(data.True))

	ns := e.GetAnonymous()
	as.Nil(ns.Declare("public-child").Bind(data.True))
	as.Nil(ns.Private("private-child").Bind(data.True))

	_, err := env.ResolveSymbol(ns, LS("public-child"))
	as.Nil(err)

	ent := mustResolveSymbol(ns, LS("private-child"))
	as.NotNil(ent)

	_, err = env.ResolveSymbol(ns, LS("public-parent"))
	as.Nil(err)

	ls := LS("private-parent")
	defer as.ExpectPanic(fmt.Errorf(env.ErrNameNotDeclared, ls))
	mustResolveSymbol(ns, ls)
}

func TestResolveValue(t *testing.T) {
	as := assert.New(t)

	e := env.NewEnvironment()
	root := e.GetRoot()
	as.Nil(root.Declare("public-parent").Bind(data.True))
	as.Nil(root.Private("private-parent").Bind(data.True))

	ns := e.GetAnonymous()
	as.Nil(ns.Declare("public-child").Bind(data.True))
	as.Nil(ns.Private("private-child").Bind(data.True))

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
