package env_test

import (
	"testing"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/env"
	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
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

	e, ok = ns1.Resolve("second-child")
	as.False(ok)
}
