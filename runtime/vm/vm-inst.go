package vm

import (
	"errors"
	"slices"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/env"
	"github.com/kode4food/ale/internal/sequence"
)

func doNull(vm *VM) {
	vm.MEM[vm.SP] = data.Null
	vm.SP--
	vm.PC++
}

func doZero(vm *VM) {
	vm.MEM[vm.SP] = data.Integer(0)
	vm.SP--
	vm.PC++
}

func doPosInt(vm *VM) {
	vm.MEM[vm.SP] = data.Integer(vm.INST.Operand())
	vm.SP--
	vm.PC++
}

func doNegInt(vm *VM) {
	vm.MEM[vm.SP] = -data.Integer(vm.INST.Operand())
	vm.SP--
	vm.PC++
}

func doTrue(vm *VM) {
	vm.MEM[vm.SP] = data.True
	vm.SP--
	vm.PC++
}

func doFalse(vm *VM) {
	vm.MEM[vm.SP] = data.False
	vm.SP--
	vm.PC++
}

func doConst(vm *VM) {
	vm.MEM[vm.SP] = vm.CL.Constants[vm.INST.Operand()]
	vm.SP--
	vm.PC++
}

func doArg(vm *VM) {
	vm.MEM[vm.SP] = vm.ARGS[vm.INST.Operand()]
	vm.SP--
	vm.PC++
}

func doRestArg(vm *VM) {
	vm.MEM[vm.SP] = vm.ARGS[vm.INST.Operand():]
	vm.SP--
	vm.PC++
}

func doArgLen(vm *VM) {
	vm.MEM[vm.SP] = data.Integer(len(vm.ARGS))
	vm.SP--
	vm.PC++
}

func doClosure(vm *VM) {
	vm.MEM[vm.SP] = vm.CL.Captured[vm.INST.Operand()]
	vm.SP--
	vm.PC++
}

func doLoad(vm *VM) {
	vm.MEM[vm.SP] = vm.MEM[vm.LP+int(vm.INST.Operand())]
	vm.SP--
	vm.PC++
}

func doStore(vm *VM) {
	vm.SP++
	vm.MEM[vm.LP+int(vm.INST.Operand())] = vm.MEM[vm.SP]
	vm.PC++
}

func doNewRef(vm *VM) {
	vm.MEM[vm.SP] = new(Ref)
	vm.SP--
	vm.PC++
}

func doBindRef(vm *VM) {
	vm.SP++
	ref := vm.MEM[vm.SP].(*Ref)
	vm.SP++
	ref.Value = vm.MEM[vm.SP]
	vm.PC++
}

func doDeref(vm *VM) {
	SP1 := vm.SP + 1
	vm.MEM[SP1] = vm.MEM[SP1].(*Ref).Value
	vm.PC++
}

func doCar(vm *VM) {
	SP1 := vm.SP + 1
	vm.MEM[SP1] = vm.MEM[SP1].(data.Pair).Car()
	vm.PC++
}

func doCdr(vm *VM) {
	SP1 := vm.SP + 1
	vm.MEM[SP1] = vm.MEM[SP1].(data.Pair).Cdr()
	vm.PC++
}

func doCons(vm *VM) {
	vm.SP++
	SP1 := vm.SP + 1
	if p, ok := vm.MEM[SP1].(data.Prepender); ok {
		vm.MEM[SP1] = p.Prepend(vm.MEM[vm.SP])
		vm.PC++
		return
	}
	vm.MEM[SP1] = data.NewCons(vm.MEM[vm.SP], vm.MEM[SP1])
	vm.PC++
}

func doEmpty(vm *VM) {
	SP1 := vm.SP + 1
	vm.MEM[SP1] = data.Bool(vm.MEM[SP1].(data.Sequence).IsEmpty())
	vm.PC++
}

func doEq(vm *VM) {
	vm.SP++
	SP1 := vm.SP + 1
	vm.MEM[SP1] = data.Bool(vm.MEM[SP1].Equal(vm.MEM[vm.SP]))
	vm.PC++
}

func doNot(vm *VM) {
	SP1 := vm.SP + 1
	vm.MEM[SP1] = !vm.MEM[SP1].(data.Bool)
	vm.PC++
}

func doDeclare(vm *VM) {
	vm.SP++
	vm.CL.Globals.Declare(
		vm.MEM[vm.SP].(data.Local),
	)
	vm.PC++
}

func doPrivate(vm *VM) {
	vm.SP++
	vm.CL.Globals.Private(
		vm.MEM[vm.SP].(data.Local),
	)
	vm.PC++
}

func doBind(vm *VM) {
	vm.SP++
	name := vm.MEM[vm.SP].(data.Local)
	vm.SP++
	vm.CL.Globals.Declare(name).Bind(vm.MEM[vm.SP])
	vm.PC++
}

func doResolve(vm *VM) {
	SP1 := vm.SP + 1
	vm.MEM[SP1] = env.MustResolveValue(
		vm.CL.Globals,
		vm.MEM[SP1].(data.Symbol),
	)
	vm.PC++
}

func doDup(vm *VM) {
	vm.MEM[vm.SP] = vm.MEM[vm.SP+1]
	vm.SP--
	vm.PC++
}

func doPop(vm *VM) {
	vm.SP++
	vm.PC++
}

func doAdd(vm *VM) {
	vm.SP++
	SP1 := vm.SP + 1
	vm.MEM[SP1] = vm.MEM[SP1].(data.Number).Add(
		vm.MEM[vm.SP].(data.Number),
	)
	vm.PC++
}

func doSub(vm *VM) {
	vm.SP++
	SP1 := vm.SP + 1
	vm.MEM[SP1] = vm.MEM[SP1].(data.Number).Sub(
		vm.MEM[vm.SP].(data.Number),
	)
	vm.PC++
}

func doMul(vm *VM) {
	vm.SP++
	SP1 := vm.SP + 1
	vm.MEM[SP1] = vm.MEM[SP1].(data.Number).Mul(
		vm.MEM[vm.SP].(data.Number),
	)
	vm.PC++
}

func doDiv(vm *VM) {
	vm.SP++
	SP1 := vm.SP + 1
	vm.MEM[SP1] = vm.MEM[SP1].(data.Number).Div(
		vm.MEM[vm.SP].(data.Number),
	)
	vm.PC++
}

func doMod(vm *VM) {
	vm.SP++
	SP1 := vm.SP + 1
	vm.MEM[SP1] = vm.MEM[SP1].(data.Number).Mod(
		vm.MEM[vm.SP].(data.Number),
	)
	vm.PC++
}

func doNumEq(vm *VM) {
	vm.SP++
	SP1 := vm.SP + 1
	vm.MEM[SP1] = data.Bool(
		data.EqualTo == vm.MEM[SP1].(data.Number).Cmp(
			vm.MEM[vm.SP].(data.Number),
		),
	)
	vm.PC++
}

func doNumLt(vm *VM) {
	vm.SP++
	SP1 := vm.SP + 1
	vm.MEM[SP1] = data.Bool(
		data.LessThan == vm.MEM[SP1].(data.Number).Cmp(
			vm.MEM[vm.SP].(data.Number),
		),
	)
	vm.PC++
}

func doNumLte(vm *VM) {
	vm.SP++
	SP1 := vm.SP + 1
	cmp := vm.MEM[SP1].(data.Number).Cmp(
		vm.MEM[vm.SP].(data.Number),
	)
	vm.MEM[SP1] = data.Bool(
		cmp == data.LessThan || cmp == data.EqualTo,
	)
	vm.PC++
}

func doNumGt(vm *VM) {
	vm.SP++
	SP1 := vm.SP + 1
	vm.MEM[SP1] = data.Bool(
		data.GreaterThan == vm.MEM[SP1].(data.Number).Cmp(
			vm.MEM[vm.SP].(data.Number),
		),
	)
	vm.PC++
}

func doNumGte(vm *VM) {
	vm.SP++
	SP1 := vm.SP + 1
	cmp := vm.MEM[SP1].(data.Number).Cmp(
		vm.MEM[vm.SP].(data.Number),
	)
	vm.MEM[SP1] = data.Bool(
		cmp == data.GreaterThan || cmp == data.EqualTo,
	)
	vm.PC++
}

func doNeg(vm *VM) {
	SP1 := vm.SP + 1
	vm.MEM[SP1] = data.Integer(0).Sub(
		vm.MEM[SP1].(data.Number),
	)
	vm.PC++
}

func doCall0(vm *VM) {
	SP1 := vm.SP + 1
	vm.MEM[SP1] = vm.MEM[SP1].(data.Procedure).Call()
	vm.PC++
}

func doCall1(vm *VM) {
	vm.SP++
	SP1 := vm.SP + 1
	vm.MEM[SP1] = vm.MEM[vm.SP].(data.Procedure).Call(vm.MEM[SP1])
	vm.PC++
}

func doCall(vm *VM) {
	op := vm.INST.Operand()
	SP1 := vm.SP + 1
	SP2 := SP1 + 1
	fn := vm.MEM[SP1].(data.Procedure)
	args := slices.Clone(vm.MEM[SP2 : SP2+int(op)])
	RES := SP1 + int(op)
	vm.MEM[RES] = fn.Call(args...)
	vm.SP = RES - 1
	vm.PC++
}

func doCallWith(vm *VM) {
	vm.SP++
	SP1 := vm.SP + 1
	vm.MEM[SP1] = vm.MEM[vm.SP].(data.Procedure).Call(
		sequence.ToValues(vm.MEM[SP1].(data.Sequence))...,
	)
	vm.PC++
}

func doTailCall(vm *VM) {
	op := vm.INST.Operand()
	SP1 := vm.SP + 1
	SP2 := SP1 + 1
	val := vm.MEM[SP1]
	vm.ARGS = slices.Clone(vm.MEM[SP2 : SP2+int(op)])
	cl, ok := val.(*Closure)
	if !ok {
		vm.ST = SUCCESS
		vm.RES = val.(data.Procedure).Call(vm.ARGS...)
		return
	}
	if cl == vm.CL {
		vm.initState()
		return
	}
	vm.CL = cl
	if len(vm.MEM) < vm.CL.StackSize+vm.CL.LocalCount {
		vm.initMem()
		return
	}
	vm.initCode()
}

func doJump(vm *VM) {
	vm.PC = int(vm.INST.Operand())
}

func doCondJump(vm *VM) {
	vm.SP++
	if vm.MEM[vm.SP] != data.False {
		vm.PC = int(vm.INST.Operand())
		return
	}
	vm.PC++
}

func doNoOp(vm *VM) {
	vm.PC++
}

func doPanic(vm *VM) {
	vm.ST = FAILURE
	panic(errors.New(data.ToString(vm.MEM[vm.SP+1])))
}

func doReturn(vm *VM) {
	vm.ST = SUCCESS
	vm.RES = vm.MEM[vm.SP+1]
}

func doRetNull(vm *VM) {
	vm.ST = SUCCESS
	vm.RES = data.Null
}

func doRetTrue(vm *VM) {
	vm.ST = SUCCESS
	vm.RES = data.True
}

func doRetFalse(vm *VM) {
	vm.ST = SUCCESS
	vm.RES = data.False
}
