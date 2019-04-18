package builtin_test

import (
	"testing"

	"gitlab.com/kode4food/ale/api"
	"gitlab.com/kode4food/ale/bootstrap"
	"gitlab.com/kode4food/ale/bootstrap/builtin"

	"gitlab.com/kode4food/ale/internal/assert"
)

func TestPredicates(t *testing.T) {
	as := assert.New(t)

	manager := bootstrap.DevNullManager()
	bootstrap.Into(manager)

	f1 := api.ApplicativeFunction(builtin.Str)
	as.False(builtin.IsSpecial(f1))
	as.True(builtin.IsApply(f1))

	ifFunc, ok := manager.GetRoot().Resolve("if")
	as.True(ok)
	as.True(builtin.IsSpecial(ifFunc))
	as.False(builtin.IsApply(ifFunc))
}
