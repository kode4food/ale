package builtin_test

import (
	"testing"

	"gitlab.com/kode4food/ale/bootstrap/builtin"
	"gitlab.com/kode4food/ale/data"
	"gitlab.com/kode4food/ale/internal/assert"
	. "gitlab.com/kode4food/ale/internal/assert/helpers"
)

func TestRead(t *testing.T) {
	as := assert.New(t)

	r1 := builtin.Read(S("[1 2 3]")).(data.Vector)

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
		data.Call(func(_ ...data.Value) data.Value {
			builtin.Raise(S("blowed up!"))
			return S("wrong")
		}),
		data.Call(func(args ...data.Value) data.Value {
			as.String("blowed up!", args[0])
			triggered = true
			return data.Nil
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
		data.Call(func(_ ...data.Value) data.Value {
			builtin.Raise(S("blowed up!"))
			return S("wrong")
		}),
		data.Call(func(_ ...data.Value) data.Value {
			triggered = true
			return data.Nil
		}),
	)
}

func TestDoEval(t *testing.T) {
	as := assert.New(t)
	as.EvalTo(`
		(do
			55
			(if true 99 33))
	`, F(99))
}

func TestTrueFalseEval(t *testing.T) {
	as := assert.New(t)
	as.EvalTo(`true`, data.True)
	as.EvalTo(`false`, data.False)
	as.EvalTo(`nil`, data.Nil)
}

func TestReadEval(t *testing.T) {
	as := assert.New(t)
	as.EvalTo(`
		(eval (read "(str \"hello\" \"you\" \"test\")"))
	`, S("helloyoutest"))
}
