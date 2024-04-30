package vm

import "github.com/kode4food/ale/pkg/data"

func doConst(vm *VM) {
	vm.MEM[vm.SP] = vm.CL.Constants[vm.INST.Operand()]
	vm.SP--
	vm.PC++
}

func doFalse(vm *VM) {
	vm.MEM[vm.SP] = data.False
	vm.SP--
	vm.PC++
}

func doNegInt(vm *VM) {
	vm.MEM[vm.SP] = -data.Integer(vm.INST.Operand())
	vm.SP--
	vm.PC++
}

func doNull(vm *VM) {
	vm.MEM[vm.SP] = data.Null
	vm.SP--
	vm.PC++
}

func doPosInt(vm *VM) {
	vm.MEM[vm.SP] = data.Integer(vm.INST.Operand())
	vm.SP--
	vm.PC++
}

func doTrue(vm *VM) {
	vm.MEM[vm.SP] = data.True
	vm.SP--
	vm.PC++
}

func doZero(vm *VM) {
	vm.MEM[vm.SP] = data.Integer(0)
	vm.SP--
	vm.PC++
}
