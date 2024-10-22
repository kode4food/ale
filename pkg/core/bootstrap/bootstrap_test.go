package bootstrap_test

import (
	"testing"

	"github.com/kode4food/ale/internal/assert"
	"github.com/kode4food/ale/internal/compiler/special"
	"github.com/kode4food/ale/pkg/core/bootstrap"
	"github.com/kode4food/ale/pkg/data"
	"github.com/kode4food/ale/pkg/env"
)

func TestDevNullEnvironment(t *testing.T) {
	as := assert.New(t)

	e := bootstrap.DevNullEnvironment()
	ns := e.GetRoot()

	// It's okay to snapshot an environment if nobody has attempted to resolve
	// an unbound namespace value
	as.IsNotBound(ns, "*args*")
	v, ok := as.IsBound(ns, "*in*").(data.Sequence)
	as.True(ok)
	as.True(v.IsEmpty())
}

func TestTopLevelEnvironment(t *testing.T) {
	as := assert.New(t)

	e := bootstrap.TopLevelEnvironment()
	ns := e.GetRoot()

	_, ok := as.IsBound(ns, "*args*").(data.Vector)
	as.True(ok)
}

func TestBootstrapInto(t *testing.T) {
	as := assert.New(t)

	e := bootstrap.TopLevelEnvironment()
	ns := e.GetRoot()

	_, ok := as.IsBound(ns, "define*").(special.Call)
	as.True(ok)
}

func BenchmarkBootstrapping(b *testing.B) {
	for n := 0; n < b.N; n++ {
		e := env.NewEnvironment()
		bootstrap.DevNull(e)
		bootstrap.Into(e)
	}
}
