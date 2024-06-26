package env_test

import (
	"fmt"
	"testing"

	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/pkg/data"
	"github.com/kode4food/ale/pkg/env"
)

func TestSnapshot(t *testing.T) {
	as := assert.New(t)

	e1 := env.NewEnvironment()
	root := e1.GetRoot()
	root.Declare("public-parent").Bind(data.True)
	root.Private("private-parent").Bind(data.True)

	ns1 := e1.GetQualified("some-ns")
	ns1.Declare("public-child").Bind(data.True)
	ns1.Private("private-child").Bind(data.True)

	e2 := env.NewEnvironment()
	ns2, err := ns1.Snapshot(e2)
	as.Nil(err)
	as.Equal(LS("some-ns"), ns2.Domain())
	as.Equal(e2, ns2.Environment())

	ns2.Declare("second-child").Bind(data.True)
	as.NotNil(ns2)
	as.Nil(err)

	d := ns2.Declared()
	as.Equal(2, len(d))
	e, ok := ns2.Resolve("public-child")
	as.True(ok)
	as.Equal(data.True, e.Value())

	e, ok = ns2.Resolve("second-child")
	as.True(ok)
	as.Equal(data.True, e.Value())

	_, ok = ns1.Resolve("second-child")
	as.False(ok)
}

func TestChainedSnapshotErrors(t *testing.T) {
	as := assert.New(t)

	e1 := env.NewEnvironment()
	root := e1.GetRoot()
	ns1 := e1.GetQualified("some-ns")

	sym1 := data.Local("was-unbound-but-resolved")
	ns1.Declare(sym1)
	e, ok := ns1.Resolve(sym1)
	as.True(ok)

	e2, err := e1.Snapshot()
	as.Nil(e2)
	as.EqualError(err, fmt.Sprintf(env.ErrSnapshotIncomplete, sym1))

	e.Bind(data.True)
	e2, err = e1.Snapshot()
	as.NotNil(e2)
	as.Nil(err)

	sym2 := data.Local("also-unbound-but-resolved")
	root.Declare(sym2)
	_, ok = root.Resolve(sym2)
	as.True(ok)

	_, err = ns1.Snapshot(env.NewEnvironment())
	as.EqualError(err, fmt.Sprintf(env.ErrSnapshotIncomplete, sym2))
}
