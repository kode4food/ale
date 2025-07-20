package env_test

import (
	"testing"

	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/env"
)

func TestSnapshot(t *testing.T) {
	as := assert.New(t)

	e1 := env.NewEnvironment()
	root := e1.GetRoot()
	as.NoError(env.BindPublic(root, "public-parent", data.True))
	as.NoError(env.BindPrivate(root, "private-parent", data.True))

	ns1 := env.MustGetQualified(e1, "some-ns")
	as.NoError(env.BindPublic(ns1, "public-child", data.True))
	as.NoError(env.BindPrivate(ns1, "private-child", data.True))

	e2 := env.NewEnvironment()
	ns2 := ns1.Snapshot(e2)
	as.Equal(LS("some-ns"), ns2.Domain())
	as.Equal(e2, ns2.Environment())

	as.NoError(env.BindPublic(ns2, "second-child", data.True))
	as.NotNil(ns2)

	d := ns2.Declared()
	as.Equal(2, len(d))
	as.Equal(data.True, as.IsBound(ns2, "public-child"))
	as.Equal(data.True, as.IsBound(ns2, "second-child"))
	as.IsNotDeclared(ns1, "second-child")
}

func TestChainedSnapshotErrors(t *testing.T) {
	as := assert.New(t)

	e1 := env.NewEnvironment()
	root := e1.GetRoot()
	ns1 := env.MustGetQualified(e1, "some-ns")

	sym1 := data.Local("was-unbound-but-resolved")
	e, err := ns1.Public(sym1)
	if as.NoError(err) {
		as.IsNotBound(ns1, sym1)
	}

	as.NoError(e.Bind(data.True))
	e2 := e1.Snapshot()
	as.NotNil(e2)

	sym2 := data.Local("also-unbound-but-resolved")
	_, err = root.Public(sym2)
	if as.NoError(err) {
		as.IsNotBound(root, sym2)
	}
}
