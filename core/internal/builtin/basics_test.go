package builtin_test

import (
	"testing"

	"github.com/kode4food/ale/core/internal/builtin"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
)

func TestRead(t *testing.T) {
	as := assert.New(t)

	r1 := builtin.Read.Call(S("[1 2 3]")).(data.Vector)

	v2, ok := r1.ElementAt(0)
	as.True(ok)
	as.Number(1, v2)

	v3, ok := r1.ElementAt(2)
	as.True(ok)
	as.Number(3, v3)
}

func TestEmptyRead(t *testing.T) {
	as := assert.New(t)
	r1 := builtin.Read.Call(S(""))
	as.Nil(r1)
}

func TestRaise(t *testing.T) {
	as := assert.New(t)

	defer func() {
		if rec := recover(); rec != nil {
			as.String("blowed up!", rec)
			return
		}
		as.Fail("proper error not raised")
	}()

	builtin.Raise.Call(S("blowed up!"))
}

func TestRecover(t *testing.T) {
	as := assert.New(t)
	var triggered = false
	builtin.Recover.Call(
		data.Applicative(func(...data.Value) data.Value {
			builtin.Raise.Call(S("blowed up!"))
			return S("wrong")
		}, 0),
		data.Applicative(func(args ...data.Value) data.Value {
			as.String("blowed up!", args[0])
			triggered = true
			return data.Nil
		}, 1),
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

	builtin.Defer.Call(
		data.Applicative(func(...data.Value) data.Value {
			builtin.Raise.Call(S("blowed up!"))
			return S("wrong")
		}, 0),
		data.Applicative(func(...data.Value) data.Value {
			triggered = true
			return data.Nil
		}, 0),
	)
}

func TestBeginEval(t *testing.T) {
	as := assert.New(t)
	as.EvalTo(`
		(begin
			55
			(if true 99 33))
	`, F(99))
}

func TestTrueFalseEval(t *testing.T) {
	as := assert.New(t)
	as.EvalTo(`true`, data.True)
	as.EvalTo(`false`, data.False)
	as.EvalTo(`'()`, data.Nil)
}

func TestReadEval(t *testing.T) {
	as := assert.New(t)
	as.EvalTo(`
		(eval (read "(str \"hello\" \"you\" \"test\")"))
	`, S("helloyoutest"))
}
