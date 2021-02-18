package bootstrap_test

import (
	"testing"

	"github.com/kode4food/ale/compiler/encoder"
	"github.com/kode4food/ale/core/bootstrap"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/assert"
	"github.com/kode4food/ale/internal/stream"
)

func TestDevNullEnvironment(t *testing.T) {
	as := assert.New(t)

	e := bootstrap.DevNullEnvironment()
	ns := e.GetRoot()

	_, ok := ns.Resolve("*args*")
	as.False(ok)

	i, ok := ns.Resolve("*in*")
	as.True(ok && i.IsBound())
	r, ok := i.Value().(stream.Reader)
	as.True(ok)
	as.True(r.IsEmpty())
}

func TestTopLevelEnvironment(t *testing.T) {
	as := assert.New(t)

	e := bootstrap.TopLevelEnvironment()
	ns := e.GetRoot()

	a, ok := ns.Resolve("*args*")
	as.True(ok && a.IsBound())

	_, ok = a.Value().(data.Vector)
	as.True(ok)
}

func TestBootstrapInto(t *testing.T) {
	as := assert.New(t)

	e := bootstrap.TopLevelEnvironment()
	bootstrap.Into(e)
	ns := e.GetRoot()

	d, ok := ns.Resolve("define*")
	as.True(ok && d.IsBound())

	_, ok = d.Value().(encoder.Call)
	as.True(ok)
}

func BenchmarkBootstrapping(b *testing.B) {
	for n := 0; n < b.N; n++ {
		e := bootstrap.DevNullEnvironment()
		bootstrap.Into(e)
	}
}
