package builtin_test

import (
	"fmt"
	"testing"

	"github.com/kode4food/ale/core/builtin"
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

func testRecover(as *assert.Wrapper, err any, errStr string) {
	var triggered = false
	builtin.Recover.Call(
		data.MakeProcedure(func(...data.Value) data.Value {
			panic(err)
		}, 0),
		data.MakeProcedure(func(args ...data.Value) data.Value {
			as.String(errStr, args[0])
			triggered = true
			return data.Null
		}, 1),
	)
	as.True(triggered)
}

func TestRecover(t *testing.T) {
	as := assert.New(t)

	errStr := "blew up"
	testRecover(as, S(errStr), errStr)
	testRecover(as, fmt.Errorf(errStr), errStr)

	defer as.ExpectProgrammerError("recover returned invalid result")
	testRecover(as, &struct{}{}, "won't be needed")
}

func TestDefer(t *testing.T) {
	as := assert.New(t)
	var triggered = false

	defer func() {
		as.True(triggered)
		recover()
	}()

	builtin.Defer.Call(
		data.MakeProcedure(func(...data.Value) data.Value {
			panic(S("blowed up!"))
		}, 0),
		data.MakeProcedure(func(...data.Value) data.Value {
			triggered = true
			return data.Null
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
	as.EvalTo(`'()`, data.Null)
}

func TestReadEval(t *testing.T) {
	as := assert.New(t)
	as.EvalTo(`
		(eval (read "(str \"hello\" \"you\" \"test\")"))
	`, S("helloyoutest"))
}
