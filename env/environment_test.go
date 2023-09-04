package env_test

import (
	"fmt"
	"testing"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/env"
	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
)

func TestResolveSymbol(t *testing.T) {
	as := assert.New(t)

	e := env.NewEnvironment()
	root := e.GetRoot()
	root.Declare("public-parent").Bind(data.True)
	root.Private("private-parent").Bind(data.True)

	ns := e.GetAnonymous()
	ns.Declare("public-child").Bind(data.True)
	ns.Private("private-child").Bind(data.True)

	_, ok := env.ResolveSymbol(ns, LS("public-child"))
	as.True(ok)

	ent := env.MustResolveSymbol(ns, LS("private-child"))
	as.NotNil(ent)

	_, ok = env.ResolveSymbol(ns, LS("public-parent"))
	as.True(ok)

	ls := LS("private-parent")
	defer as.ExpectPanic(fmt.Sprintf(env.ErrSymbolNotDeclared, ls))
	env.MustResolveSymbol(ns, ls)
}

func TestResolveValue(t *testing.T) {
	as := assert.New(t)

	e := env.NewEnvironment()
	root := e.GetRoot()
	root.Declare("public-parent").Bind(data.True)
	root.Private("private-parent").Bind(data.True)

	ns := e.GetAnonymous()
	ns.Declare("public-child").Bind(data.True)
	ns.Private("private-child").Bind(data.True)

	res, ok := env.ResolveValue(ns, LS("public-child"))
	as.True(res)
	as.True(ok)

	res = env.MustResolveValue(ns, LS("private-child"))
	as.True(res)

	_, ok = env.ResolveValue(ns, LS("public-parent"))
	as.True(ok)

	ls := LS("private-parent")
	defer as.ExpectPanic(fmt.Sprintf(env.ErrSymbolNotBound, ls))
	env.MustResolveValue(ns, ls)
}
