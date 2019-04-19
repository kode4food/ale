package builtin_test

import (
	"testing"

	"gitlab.com/kode4food/ale/bootstrap"
	"gitlab.com/kode4food/ale/bootstrap/builtin"
	"gitlab.com/kode4food/ale/data"
	"gitlab.com/kode4food/ale/internal/assert"
	. "gitlab.com/kode4food/ale/internal/assert/helpers"
)

func getCall(v data.Value) data.Call {
	return v.(data.Caller).Caller()
}

func TestApply(t *testing.T) {
	as := assert.New(t)

	vCall := data.Call(builtin.Vector)
	as.True(builtin.IsApply(vCall))
	as.False(builtin.IsApply(S("55")))

	v1 := builtin.Vector(S("4"), S("5"), S("6"))
	v2 := builtin.Apply(vCall, S("1"), S("2"), S("3"), v1)
	v3 := builtin.Apply(vCall, v1)

	as.String(`["4" "5" "6"]`, v1)
	as.String(`["1" "2" "3" "4" "5" "6"]`, v2)
	as.String(`["4" "5" "6"]`, v3)
}

func TestPartial(t *testing.T) {
	as := assert.New(t)

	vCall := data.Call(builtin.Vector)
	p1 := builtin.Partial(vCall, S("1"), S("2"))
	v1 := getCall(p1)(S("3"), S("4"), S("5"))
	v2 := getCall(p1)(S("7"), S("9"))

	as.String(`["1" "2" "3" "4" "5"]`, v1)
	as.String(`["1" "2" "7" "9"]`, v2)
}

func TestPredicates(t *testing.T) {
	as := assert.New(t)

	manager := bootstrap.DevNullManager()
	bootstrap.Into(manager)

	f1 := data.ApplicativeFunction(builtin.Str)
	as.False(builtin.IsSpecial(f1))
	as.True(builtin.IsApply(f1))

	ifFunc, ok := manager.GetRoot().Resolve("if")
	as.True(ok)
	as.True(builtin.IsSpecial(ifFunc))
	as.False(builtin.IsApply(ifFunc))
}
