package vm

import (
	"slices"

	"github.com/kode4food/ale/internal/sequence"
	"github.com/kode4food/ale/pkg/data"
)

func doCall(vm *machine) {
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

func doCall0(vm *machine) {
	SP1 := vm.SP + 1
	vm.MEM[SP1] = vm.MEM[SP1].(data.Procedure).Call()
	vm.PC++
}

func doCall1(vm *machine) {
	vm.SP++
	SP1 := vm.SP + 1
	vm.MEM[SP1] = vm.MEM[vm.SP].(data.Procedure).Call(vm.MEM[SP1])
	vm.PC++
}

func doCallWith(vm *machine) {
	vm.SP++
	SP1 := vm.SP + 1
	vm.MEM[SP1] = vm.MEM[vm.SP].(data.Procedure).Call(
		sequence.ToValues(vm.MEM[SP1].(data.Sequence))...,
	)
	vm.PC++
}

func doTailCall(vm *machine) {
	op := vm.INST.Operand()
	SP1 := vm.SP + 1
	SP2 := SP1 + 1
	val := vm.MEM[SP1]
	vm.ARGS = slices.Clone(vm.MEM[SP2 : SP2+int(op)])
	cl, ok := val.(*Closure)
	if !ok {
		vm.ST = success
		vm.RES = val.(data.Procedure).Call(vm.ARGS...)
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
