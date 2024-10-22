package assert_test

import (
	"testing"

	"github.com/kode4food/ale/internal/assert"
	"github.com/kode4food/ale/pkg/data"
)

func TestIsBound(t *testing.T) {
	as := assert.New(t)
	ns := assert.GetTestNamespace()
	as.Nil(ns.Declare("found").Bind(data.True))
	as.True(as.IsBound(ns, "found"))
}

func TestIsNotBound(t *testing.T) {
	as := assert.New(t)
	ns := assert.GetTestNamespace()
	as.NotNil(ns.Declare("not-bound"))
	as.IsNotBound(ns, "not-bound")
}

func TestIsNotDeclared(t *testing.T) {
	as := assert.New(t)
	ns := assert.GetTestNamespace()
	as.IsNotDeclared(ns, "not-declared")
}
