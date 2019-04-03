package builtin_test

import (
	"testing"

	"gitlab.com/kode4food/ale/api"
	"gitlab.com/kode4food/ale/internal/assert"
	. "gitlab.com/kode4food/ale/internal/assert/helpers"
	"gitlab.com/kode4food/ale/internal/builtin"
)

func TestRead(t *testing.T) {
	as := assert.New(t)

	r1 := builtin.Read(S("[1 2 3]")).(api.Vector)

	v2, ok := r1.ElementAt(0)
	as.True(ok)
	as.Integer(1, v2)

	v3, ok := r1.ElementAt(2)
	as.True(ok)
	as.Integer(3, v3)
}
