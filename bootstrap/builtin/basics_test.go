package builtin_test

import (
	"testing"

	"gitlab.com/kode4food/ale/api"
	"gitlab.com/kode4food/ale/bootstrap/builtin"
	"gitlab.com/kode4food/ale/internal/assert"
	. "gitlab.com/kode4food/ale/internal/assert/helpers"
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

func TestEmptyRead(t *testing.T) {
	as := assert.New(t)
	r1 := builtin.Read(S(""))
	as.Nil(r1)
}

func TestRaise(t *testing.T) {
	as := assert.New(t)

	defer func() {
		if rec := recover(); rec != nil {
			as.String("blowed up!", rec)
			return
		}
		as.Fail("error not raised")
	}()

	builtin.Raise(S("blowed up!"))
}

func TestRecover(t *testing.T) {
	as := assert.New(t)
	var triggered = false
	builtin.Recover(
		api.Call(func(_ ...api.Value) api.Value {
			builtin.Raise(S("blowed up!"))
			return S("wrong")
		}),
		api.Call(func(args ...api.Value) api.Value {
			as.String("blowed up!", args[0])
			triggered = true
			return api.Nil
		}),
	)
	as.True(triggered)
}

func TestDefer(t *testing.T) {
	as := assert.New(t)
	var triggered = false

	defer func() {
		as.True(triggered)
		recover()
	}()

	builtin.Defer(
		api.Call(func(_ ...api.Value) api.Value {
			builtin.Raise(S("blowed up!"))
			return S("wrong")
		}),
		api.Call(func(_ ...api.Value) api.Value {
			triggered = true
			return api.Nil
		}),
	)
}
