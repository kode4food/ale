package vm

import (
	"slices"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/sequence"
)

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
	vm.ARGS.Data = slices.Clone(vm.MEM[SP2 : SP2+int(op)])
	cl, ok := val.(*Closure)
	if !ok {
		vm.ST = SUCCESS
		vm.RES = val.(data.Procedure).Call(vm.ARGS.Data...)
		return
	}
	if cl == vm.CL {
		vm.initState()
		return
	}
	vm.CL = cl
	if len(vm.MEM) < int(vm.CL.StackSize+vm.CL.LocalCount) {
		vm.initMem()
		return
	}
	vm.initCode()
}
