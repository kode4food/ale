package builtin_test

import (
	"testing"

	"gitlab.com/kode4food/ale/api"
	"gitlab.com/kode4food/ale/internal/assert"
	. "gitlab.com/kode4food/ale/internal/assert/helpers"
	"gitlab.com/kode4food/ale/internal/builtin"
)

func TestAssoc(t *testing.T) {
	as := assert.New(t)

	a1 := builtin.Assoc(K("hello"), S("foo"))
	m1 := a1.(api.Mapped)
	v1, ok := m1.Get(K("hello"))
	as.True(ok)
	as.String("foo", v1)

	as.True(builtin.IsAssoc(a1))
	as.False(builtin.IsAssoc(I(99)))

	as.True(builtin.IsMapped(a1))
	as.False(builtin.IsMapped(I(99)))
}
