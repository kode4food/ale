package assert_test

import (
	"testing"

	"github.com/kode4food/ale/internal/assert"
	"github.com/kode4food/ale/pkg/data"
	"github.com/kode4food/ale/pkg/env"
)

func TestIsBound(t *testing.T) {
	as := assert.New(t)
	ns := assert.GetTestNamespace()
	as.Nil(env.BindPublic(ns, "found", data.True))
	as.True(as.IsBound(ns, "found"))
}

func TestIsNotBound(t *testing.T) {
	as := assert.New(t)
	ns := assert.GetTestNamespace()
	e, err := ns.Public("not-bound")
	as.NotNil(e)
	as.Nil(err)
	as.IsNotBound(ns, "not-bound")
}

func TestIsNotDeclared(t *testing.T) {
	as := assert.New(t)
	ns := assert.GetTestNamespace()
	as.IsNotDeclared(ns, "not-declared")
}
