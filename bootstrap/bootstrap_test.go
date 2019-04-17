package bootstrap_test

import (
	"testing"

	"gitlab.com/kode4food/ale/api"
	"gitlab.com/kode4food/ale/bootstrap"
	"gitlab.com/kode4food/ale/compiler/encoder"
	"gitlab.com/kode4food/ale/internal/assert"
	"gitlab.com/kode4food/ale/stdlib"
)

func TestDevNullManager(t *testing.T) {
	as := assert.New(t)

	manager := bootstrap.DevNullManager()
	ns := manager.GetRoot()

	_, ok := ns.Resolve("*args*")
	as.False(ok)

	in, ok := ns.Resolve("*in*")
	as.True(ok)
	r, ok := in.(stdlib.Reader)
	as.True(ok)
	as.False(r.IsSequence())
}

func TestTopLevelManager(t *testing.T) {
	as := assert.New(t)

	manager := bootstrap.TopLevelManager()
	ns := manager.GetRoot()

	args, ok := ns.Resolve("*args*")
	as.True(ok)

	_, ok = args.(api.Vector)
	as.True(ok)
}

func TestBootstrapInto(t *testing.T) {
	as := assert.New(t)

	manager := bootstrap.TopLevelManager()
	bootstrap.Into(manager)
	ns := manager.GetRoot()

	v, ok := ns.Resolve("def")
	as.True(ok)

	_, ok = v.(encoder.Call)
	as.True(ok)
}
