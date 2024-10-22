package env_test

import (
	"fmt"
	"testing"

	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/pkg/data"
	"github.com/kode4food/ale/pkg/env"
)

func TestDeclarations(t *testing.T) {
	as := assert.New(t)

	e := env.NewEnvironment()
	root := e.GetRoot()
	as.Equal(e, root.Environment())
	as.Equal(env.RootDomain, root.Domain())

	as.Nil(root.Declare("public2").Bind(data.True))
	as.Nil(root.Private("private").Bind(data.True))
	as.Nil(root.Declare("public1").Bind(data.True))

	n := root.Declared()
	as.Equal(2, len(n))
	as.Equal(LS("public1"), n[0])
	as.Equal(LS("public2"), n[1])

	e1, err := root.Resolve(n[0])
	as.NotNil(e1)
	as.Nil(err)

	as.Equal(n[0], e1.Name())
	as.Equal(root, e1.Owner())

	e2 := root.Declare(n[0])
	as.Equal(e1, e2)
}

func TestChaining(t *testing.T) {
	as := assert.New(t)

	e := env.NewEnvironment()
	root := e.GetRoot()
	as.Nil(root.Declare("in-parent").Bind(data.True))

	ns := e.GetAnonymous()
	as.Nil(ns.Declare("in-child").Bind(data.True))

	as.True(as.IsBound(ns, "in-parent"))
	as.True(as.IsBound(ns, "in-child"))
	as.True(as.IsBound(root, "in-parent"))
	as.IsNotDeclared(root, "in-child")
	s3 := env.RootSymbol("in-parent")
	v8, err := env.ResolveValue(ns, s3)
	as.True(v8)
	as.Nil(err)
}

func TestBinding(t *testing.T) {
	as := assert.New(t)

	e := env.NewEnvironment()
	root := e.GetRoot()
	d := root.Declare("some-name")

	v, err := d.Value()
	as.Nil(v)
	as.EqualError(err, fmt.Sprintf(env.ErrNameNotBound, d.Name()))

	err = d.Bind(S("some-value"))
	as.Nil(err)
	err = d.Bind(S("some-other-value"))
	as.EqualError(err, fmt.Sprintf(env.ErrNameAlreadyBound, d.Name()))

	v, err = d.Value()
	as.Nil(err)
	as.String("some-value", v)
}
