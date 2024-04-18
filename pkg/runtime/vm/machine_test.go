package vm_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/pkg/data"
	"github.com/kode4food/ale/pkg/env"
	"github.com/kode4food/ale/pkg/runtime/isa"
	"github.com/kode4food/ale/pkg/runtime/vm"
)

var constants = data.Vector{
	I(5),
	I(6),
	S("a thrown error"),
	data.MakeProcedure(numLoopSum),
	LS("a-name"),
}

func makeProcedure(code isa.Instructions) data.Procedure {
	proc := &vm.Procedure{
		Runnable: isa.Runnable{
			Code:       code,
			Constants:  constants,
			StackSize:  16,
			LocalCount: 10,
			Globals:    env.NewEnvironment().GetAnonymous(),
		},
	}
	closure := proc.Call(S("Closure"))
	return closure.(data.Procedure)
}

func runCode(code isa.Instructions) data.Value {
	fn := makeProcedure(code)
	return fn.Call(S("arg"))
}

func testResult(t *testing.T, res data.Value, code isa.Instructions) {
	as := assert.New(t)
	r := runCode(code)
	as.Equal(res, r)
}

func testPanic(t *testing.T, errStr string, code isa.Instructions) {
	as := assert.New(t)
	defer as.ExpectPanic(errors.New(errStr))
	runCode(code)
}

func TestSimple(t *testing.T) {
	testResult(t, I(11), isa.Instructions{
		isa.Const.New(0),
		isa.Const.New(1),
		isa.Add.New(),
		isa.Return.New(),
	})

	testResult(t, I(0), isa.Instructions{
		isa.Zero.New(),
		isa.Const.New(0),
		isa.Mul.New(),
		isa.Return.New(),
	})

	testResult(t, S("Closure"), isa.Instructions{
		isa.Closure.New(0),
		isa.Return.New(),
	})

	testResult(t, S("arg"), isa.Instructions{
		isa.Arg.New(0),
		isa.Return.New(),
	})
}

func TestPopAndDup(t *testing.T) {
	testResult(t, I(4), isa.Instructions{
		isa.PosInt.New(2),
		isa.Dup.New(),
		isa.NoOp.New(),
		isa.Add.New(),
		isa.Return.New(),
	})

	testResult(t, I(2), isa.Instructions{
		isa.PosInt.New(1),
		isa.PosInt.New(1),
		isa.PosInt.New(2),
		isa.PosInt.New(2),
		isa.Pop.New(),
		isa.Pop.New(),
		isa.Add.New(),
		isa.Return.New(),
	})

	testResult(t, I(6), isa.Instructions{
		isa.PosInt.New(1),
		isa.PosInt.New(2),
		isa.Add.New(),
		isa.PosInt.New(2),
		isa.Dup.New(),
		isa.Pop.New(),
		isa.Mul.New(),
		isa.Return.New(),
	})
}

func TestReturns(t *testing.T) {
	testResult(t, data.Null, isa.Instructions{
		isa.Null.New(),
		isa.Return.New(),
	})

	testResult(t, I(2), isa.Instructions{
		isa.PosInt.New(2),
		isa.Return.New(),
	})

	testResult(t, data.True, isa.Instructions{
		isa.True.New(),
		isa.Return.New(),
	})

	testResult(t, data.False, isa.Instructions{
		isa.False.New(),
		isa.Return.New(),
	})

	testResult(t, data.True, isa.Instructions{
		isa.RetTrue.New(),
	})

	testResult(t, data.False, isa.Instructions{
		isa.RetFalse.New(),
	})

	testResult(t, data.Null, isa.Instructions{
		isa.RetNull.New(),
	})
}

func TestUnary(t *testing.T) {
	testResult(t, I(-1), isa.Instructions{
		isa.PosInt.New(1),
		isa.NoOp.New(),
		isa.Neg.New(),
		isa.Return.New(),
	})

	testResult(t, data.False, isa.Instructions{
		isa.True.New(),
		isa.Not.New(),
		isa.Return.New(),
	})
}

func TestCalls(t *testing.T) {
	testResult(t, I(17), isa.Instructions{
		isa.Const.New(0),
		isa.Const.New(0),
		isa.Const.New(1),
		isa.PosInt.New(1),
		isa.Const.New(3),
		isa.Call.New(3),
		isa.Add.New(),
		isa.Return.New(),
	})

	testResult(t, I(5), isa.Instructions{
		isa.Const.New(0),
		isa.Const.New(3),
		isa.Call1.New(),
		isa.Return.New(),
	})

	testResult(t, I(0), isa.Instructions{
		isa.Const.New(3),
		isa.Call0.New(),
		isa.Return.New(),
	})
}

func TestMaths(t *testing.T) {
	testResult(t, I(3), isa.Instructions{
		isa.PosInt.New(1),
		isa.PosInt.New(2),
		isa.Add.New(), isa.Return.New(),
	})

	testResult(t, I(4), isa.Instructions{
		isa.PosInt.New(2),
		isa.PosInt.New(2),
		isa.Mul.New(), isa.Return.New(),
	})

	testResult(t, I(-2), isa.Instructions{
		isa.PosInt.New(2),
		isa.NegInt.New(1),
		isa.Mul.New(), isa.Return.New(),
	})

	testResult(t, I(-1), isa.Instructions{
		isa.PosInt.New(1),
		isa.PosInt.New(2),
		isa.Sub.New(), isa.Return.New(),
	})

	testResult(t, R(1, 2), isa.Instructions{
		isa.PosInt.New(1),
		isa.PosInt.New(2),
		isa.Div.New(), isa.Return.New(),
	})

	testResult(t, I(1), isa.Instructions{
		isa.PosInt.New(2), isa.PosInt.New(2), isa.Mul.New(),
		isa.PosInt.New(2), isa.Mul.New(),
		isa.PosInt.New(1), isa.Add.New(),
		isa.PosInt.New(2), isa.Mod.New(),
		isa.Return.New(),
	})
}

func TestRelational(t *testing.T) {
	testResult(t, data.True, isa.Instructions{
		isa.PosInt.New(1),
		isa.PosInt.New(1),
		isa.NumEq.New(), isa.Return.New(),
	})

	testResult(t, data.False, isa.Instructions{
		isa.PosInt.New(1),
		isa.PosInt.New(2),
		isa.NumEq.New(), isa.Return.New(),
	})

	testResult(t, data.True, isa.Instructions{
		isa.PosInt.New(1),
		isa.PosInt.New(2),
		isa.NumLt.New(), isa.Return.New(),
	})

	testResult(t, data.False, isa.Instructions{
		isa.PosInt.New(2),
		isa.PosInt.New(1),
		isa.NumLt.New(), isa.Return.New(),
	})

	testResult(t, data.True, isa.Instructions{
		isa.PosInt.New(1),
		isa.PosInt.New(2),
		isa.NumLte.New(), isa.Return.New(),
	})

	testResult(t, data.True, isa.Instructions{
		isa.PosInt.New(2),
		isa.PosInt.New(2),
		isa.NumLte.New(), isa.Return.New(),
	})

	testResult(t, data.False, isa.Instructions{
		isa.PosInt.New(2),
		isa.PosInt.New(1),
		isa.NumLte.New(), isa.Return.New(),
	})

	testResult(t, data.True, isa.Instructions{
		isa.PosInt.New(2),
		isa.PosInt.New(1),
		isa.NumGt.New(), isa.Return.New(),
	})

	testResult(t, data.False, isa.Instructions{
		isa.PosInt.New(1),
		isa.PosInt.New(2),
		isa.NumGt.New(), isa.Return.New(),
	})

	testResult(t, data.True, isa.Instructions{
		isa.PosInt.New(2),
		isa.PosInt.New(1),
		isa.NumGte.New(), isa.Return.New(),
	})

	testResult(t, data.True, isa.Instructions{
		isa.PosInt.New(2),
		isa.PosInt.New(2),
		isa.NumGte.New(), isa.Return.New(),
	})

	testResult(t, data.False, isa.Instructions{
		isa.PosInt.New(1),
		isa.PosInt.New(2),
		isa.NumGte.New(), isa.Return.New(),
	})
}

func TestLoadStore(t *testing.T) {
	testResult(t, I(4), isa.Instructions{
		isa.PosInt.New(2),
		isa.Store.New(0),
		isa.Load.New(0),
		isa.Load.New(0),
		isa.Mul.New(),
		isa.Return.New(),
	})
}

func TestRefs(t *testing.T) {
	testResult(t, I(-1), isa.Instructions{
		isa.PosInt.New(2),
		isa.Store.New(1),
		isa.NewRef.New(),
		isa.Store.New(2),
		isa.Load.New(1),
		isa.Load.New(2),
		isa.BindRef.New(),
		isa.PosInt.New(1),
		isa.Load.New(2),
		isa.Deref.New(),
		isa.Sub.New(),
		isa.Return.New(),
	})
}

func TestGlobals(t *testing.T) {
	testResult(t, I(3), isa.Instructions{
		isa.Const.New(4),
		isa.Declare.New(),
		isa.PosInt.New(2),
		isa.Const.New(4),
		isa.Bind.New(),
		isa.PosInt.New(1),
		isa.Const.New(4),
		isa.Resolve.New(),
		isa.Add.New(),
		isa.Return.New(),
	})
}

func TestJumps(t *testing.T) {
	testResult(t, I(4), isa.Instructions{
		isa.PosInt.New(2),
		isa.Jump.New(3),
		isa.PosInt.New(1),
		isa.PosInt.New(2),
		isa.Add.New(),
		isa.Return.New(),
	})

	testResult(t, I(4), isa.Instructions{
		isa.PosInt.New(2),
		isa.True.New(),
		isa.CondJump.New(6),
		isa.PosInt.New(1),
		isa.Add.New(),
		isa.Jump.New(8),
		isa.PosInt.New(2),
		isa.Add.New(),
		isa.Return.New(),
	})

	testResult(t, I(3), isa.Instructions{
		isa.PosInt.New(2),
		isa.False.New(),
		isa.CondJump.New(6),
		isa.PosInt.New(1),
		isa.Add.New(),
		isa.Jump.New(8),
		isa.PosInt.New(2),
		isa.Add.New(),
		isa.Return.New(),
	})
}

func TestArgs(t *testing.T) {
	as := assert.New(t)
	args := data.Vector{S("arg1"), S("arg2"), S("arg3"), S("arg4")}

	c1 := makeProcedure(isa.Instructions{
		isa.ArgLen.New(),
		isa.Return.New(),
	})
	r1 := c1.Call(args...)
	as.Equal(I(4), r1)

	c2 := makeProcedure(isa.Instructions{
		isa.Arg.New(1),
		isa.Return.New(),
	})
	r2 := c2.Call(args...)
	as.Equal(S("arg2"), r2)

	c3 := makeProcedure(isa.Instructions{
		isa.RestArg.New(2),
		isa.Return.New(),
	})
	r3 := c3.Call(args...)
	as.Equal(data.NewVector(S("arg3"), S("arg4")), r3)
}

func TestErrors(t *testing.T) {
	testPanic(t, "a thrown error", isa.Instructions{
		isa.Const.New(2),
		isa.Panic.New(),
	})
}

func TestForUnimplementedOpcodes(t *testing.T) {
	as := assert.New(t)
	for oc, effect := range isa.Effects {
		(func(oc isa.Opcode, effect *isa.Effect) {
			defer func() {
				if rec := recover(); rec != nil {
					switch oc {
					case isa.Label:
						// continue
					default:
						as.NotEqual(
							rec, fmt.Sprintf("unknown opcode: %s", oc),
						)
					}
				}
			}()
			switch {
			case effect.Ignore:
				return
			case effect.Operand == isa.Nothing:
				runCode(isa.Instructions{oc.New()})
			default:
				runCode(isa.Instructions{oc.New(isa.OperandMask)})
			}
		})(oc, effect)
	}
}

func TestBadOpcode(t *testing.T) {
	as := assert.New(t)
	defer as.ExpectProgrammerError(
		fmt.Sprintf(vm.ErrBadInstruction, isa.Instruction(isa.Label)),
	)
	runCode(isa.Instructions{isa.Instruction(isa.Label)})
}
