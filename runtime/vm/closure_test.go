package vm_test

import (
	"fmt"
	"testing"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/env"
	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/runtime/isa"
	"github.com/kode4food/ale/runtime/vm"
)

var constants = data.Values{
	I(5),
	I(6),
	S("a thrown error"),
	data.Applicative(numLoopSum),
	N("a-name"),
	LS("a-name"),
}

func makeCode(coders []isa.Coder) data.Function {
	code := make(isa.Instructions, len(coders))
	for i, c := range coders {
		code[i] = c.Instruction()
	}
	lambda := &vm.Lambda{
		Code:       code,
		Constants:  constants,
		StackSize:  16,
		LocalCount: 10,
		Globals:    env.NewEnvironment().GetAnonymous(),
	}
	closure := lambda.Call(S("closure"))
	return closure.(data.Function)
}

func runCode(coders []isa.Coder) data.Value {
	code := makeCode(coders)
	return code.Call(S("arg"))
}

func testResult(t *testing.T, res data.Value, code []isa.Coder) {
	as := assert.New(t)
	r := runCode(code)
	as.Equal(res, r)
}

func testPanic(t *testing.T, errStr string, code []isa.Coder) {
	as := assert.New(t)
	defer as.ExpectPanic(errStr)
	runCode(code)
}

func TestSimple(t *testing.T) {
	testResult(t, I(11), []isa.Coder{
		isa.New(isa.Const, 0),
		isa.New(isa.Const, 1),
		isa.Add,
		isa.Return,
	})

	testResult(t, I(0), []isa.Coder{
		isa.Zero,
		isa.New(isa.Const, 0),
		isa.Mul,
		isa.Return,
	})

	testResult(t, S("closure"), []isa.Coder{
		isa.New(isa.Closure, 0),
		isa.Return,
	})

	testResult(t, S("arg"), []isa.Coder{
		isa.New(isa.Arg, 0),
		isa.Return,
	})
}

func TestPopAndDup(t *testing.T) {
	testResult(t, I(4), []isa.Coder{isa.Two, isa.Dup, isa.Add, isa.Return})

	testResult(t, I(2), []isa.Coder{
		isa.One, isa.One, isa.Two, isa.Two,
		isa.Pop, isa.Pop,
		isa.Add,
		isa.Return,
	})

	testResult(t, I(6), []isa.Coder{
		isa.One, isa.Two, isa.Add,
		isa.Two, isa.Dup, isa.Pop, isa.Mul,
		isa.Return,
	})
}

func TestReturns(t *testing.T) {
	testResult(t, data.Nil, []isa.Coder{isa.Nil, isa.Return})
	testResult(t, I(2), []isa.Coder{isa.Two, isa.Return})
	testResult(t, data.True, []isa.Coder{isa.True, isa.Return})
	testResult(t, data.False, []isa.Coder{isa.False, isa.Return})

	testResult(t, data.True, []isa.Coder{isa.RetTrue})
	testResult(t, data.False, []isa.Coder{isa.RetFalse})
	testResult(t, data.Nil, []isa.Coder{isa.RetNil})
}

func TestUnary(t *testing.T) {
	testResult(t, I(-1), []isa.Coder{isa.One, isa.Neg, isa.Return})
	testResult(t, data.False, []isa.Coder{isa.True, isa.Not, isa.Return})
}

func TestMakeTruthy(t *testing.T) {
	testResult(t, data.True, []isa.Coder{
		isa.New(isa.Const, 3),
		isa.MakeTruthy,
		isa.Return,
	})

	testResult(t, data.True, []isa.Coder{
		isa.True,
		isa.MakeTruthy,
		isa.Return,
	})

	testResult(t, data.False, []isa.Coder{
		isa.False,
		isa.MakeTruthy,
		isa.Return,
	})

	testResult(t, data.False, []isa.Coder{
		isa.Nil,
		isa.MakeTruthy,
		isa.Return,
	})
}

func TestCalls(t *testing.T) {
	testResult(t, I(17), []isa.Coder{
		isa.New(isa.Const, 0),
		isa.New(isa.Const, 0),
		isa.New(isa.Const, 1),
		isa.One,
		isa.New(isa.Const, 3),
		isa.New(isa.Call, 3),
		isa.Add,
		isa.Return,
	})

	testResult(t, I(5), []isa.Coder{
		isa.New(isa.Const, 0),
		isa.New(isa.Const, 3),
		isa.Call1,
		isa.Return,
	})

	testResult(t, I(0), []isa.Coder{
		isa.New(isa.Const, 3),
		isa.Call0,
		isa.Return,
	})
}

func TestMaths(t *testing.T) {
	testResult(t, I(3), []isa.Coder{isa.One, isa.Two, isa.Add, isa.Return})
	testResult(t, I(4), []isa.Coder{isa.Two, isa.Two, isa.Mul, isa.Return})
	testResult(t, I(-2), []isa.Coder{isa.Two, isa.NegOne, isa.Mul, isa.Return})
	testResult(t, I(-1), []isa.Coder{isa.One, isa.Two, isa.Sub, isa.Return})
	testResult(t, R(1, 2), []isa.Coder{isa.One, isa.Two, isa.Div, isa.Return})

	testResult(t, I(1), []isa.Coder{
		isa.Two, isa.Two, isa.Mul,
		isa.Two, isa.Mul,
		isa.One, isa.Add,
		isa.Two, isa.Mod,
		isa.Return,
	})
}

func TestRelational(t *testing.T) {
	testResult(t, data.True, []isa.Coder{isa.One, isa.One, isa.Eq, isa.Return})
	testResult(t, data.False, []isa.Coder{isa.One, isa.Two, isa.Eq, isa.Return})
	testResult(t, data.True, []isa.Coder{isa.Two, isa.One, isa.Neq, isa.Return})
	testResult(t, data.False, []isa.Coder{isa.One, isa.One, isa.Neq, isa.Return})

	testResult(t, data.True, []isa.Coder{isa.One, isa.Two, isa.Lt, isa.Return})
	testResult(t, data.False, []isa.Coder{isa.Two, isa.One, isa.Lt, isa.Return})
	testResult(t, data.True, []isa.Coder{isa.One, isa.Two, isa.Lte, isa.Return})
	testResult(t, data.True, []isa.Coder{isa.Two, isa.Two, isa.Lte, isa.Return})
	testResult(t, data.False, []isa.Coder{isa.Two, isa.One, isa.Lte, isa.Return})

	testResult(t, data.True, []isa.Coder{isa.Two, isa.One, isa.Gt, isa.Return})
	testResult(t, data.False, []isa.Coder{isa.One, isa.Two, isa.Gt, isa.Return})
	testResult(t, data.True, []isa.Coder{isa.Two, isa.One, isa.Gte, isa.Return})
	testResult(t, data.True, []isa.Coder{isa.Two, isa.Two, isa.Gte, isa.Return})
	testResult(t, data.False, []isa.Coder{isa.One, isa.Two, isa.Gte, isa.Return})
}

func TestLoadStore(t *testing.T) {
	testResult(t, I(4), []isa.Coder{
		isa.Two,
		isa.New(isa.Store, 0),
		isa.New(isa.Load, 0),
		isa.New(isa.Load, 0),
		isa.Mul,
		isa.Return,
	})
}

func TestRefs(t *testing.T) {
	testResult(t, I(-1), []isa.Coder{
		isa.Two,
		isa.New(isa.Store, 1),
		isa.NewRef,
		isa.New(isa.Store, 2),
		isa.New(isa.Load, 1),
		isa.New(isa.Load, 2),
		isa.BindRef,
		isa.One,
		isa.New(isa.Load, 2),
		isa.Deref,
		isa.Sub,
		isa.Return,
	})
}

func TestGlobals(t *testing.T) {
	testResult(t, I(3), []isa.Coder{
		isa.New(isa.Const, 4),
		isa.Declare,
		isa.Two,
		isa.New(isa.Const, 4),
		isa.Bind,
		isa.One,
		isa.New(isa.Const, 5),
		isa.Resolve,
		isa.Add,
		isa.Return,
	})
}

func TestJumps(t *testing.T) {
	testResult(t, I(4), []isa.Coder{
		isa.Two,
		isa.New(isa.Jump, 3),
		isa.One,
		isa.Two,
		isa.Add,
		isa.Return,
	})

	testResult(t, I(4), []isa.Coder{
		isa.Two,
		isa.True,
		isa.New(isa.CondJump, 6),
		isa.One,
		isa.Add,
		isa.New(isa.Jump, 8),
		isa.Two,
		isa.Add,
		isa.Return,
	})

	testResult(t, I(3), []isa.Coder{
		isa.Two,
		isa.False,
		isa.New(isa.CondJump, 6),
		isa.One,
		isa.Add,
		isa.New(isa.Jump, 8),
		isa.Two,
		isa.Add,
		isa.Return,
	})
}

func TestArgs(t *testing.T) {
	as := assert.New(t)
	args := data.Values{S("arg1"), S("arg2"), S("arg3"), S("arg4")}

	c1 := makeCode([]isa.Coder{
		isa.ArgLen,
		isa.Return,
	})
	r1 := c1.Call(args...)
	as.Equal(I(4), r1)

	c2 := makeCode([]isa.Coder{
		isa.New(isa.Arg, 1),
		isa.Return,
	})
	r2 := c2.Call(args...)
	as.Equal(S("arg2"), r2)

	c3 := makeCode([]isa.Coder{
		isa.New(isa.RestArg, 2),
		isa.Return,
	})
	r3 := c3.Call(args...)
	as.Equal(data.NewVector(S("arg3"), S("arg4")), r3)
}

func TestErrors(t *testing.T) {
	testPanic(t, "a thrown error", []isa.Coder{
		isa.New(isa.Const, 2),
		isa.Panic,
	})
}

func TestExplosions(t *testing.T) {
	testPanic(t, "runtime error: index out of range", []isa.Coder{
		isa.Return,
	})
}

func TestBadOpcode(t *testing.T) {
	as := assert.New(t)
	badOpcode := isa.Opcode(isa.MaxWord)
	defer as.ExpectProgrammerError(
		fmt.Sprintf(
			"opcode can't be encoded as instruction: %s", badOpcode,
		),
	)
	runCode([]isa.Coder{badOpcode})
}
