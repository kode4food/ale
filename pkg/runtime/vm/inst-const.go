package vm

import "github.com/kode4food/ale/pkg/data"

func doConst(vm *machine) {
	vm.MEM[vm.SP] = vm.CL.Constants[vm.INST.Operand()]
	vm.SP--
	vm.PC++
}

func doFalse(vm *machine) {
	vm.MEM[vm.SP] = data.False
	vm.SP--
	vm.PC++
}

func doNegInt(vm *machine) {
	vm.MEM[vm.SP] = -data.Integer(vm.INST.Operand())
	vm.SP--
	vm.PC++
}

func doNull(vm *machine) {
	vm.MEM[vm.SP] = data.Null
	vm.SP--
	vm.PC++
}

func doPosInt(vm *machine) {
	vm.MEM[vm.SP] = data.Integer(vm.INST.Operand())
	vm.SP--
	vm.PC++
}

func doTrue(vm *machine) {
	vm.MEM[vm.SP] = data.True
	vm.SP--
	vm.PC++
}

func doZero(vm *machine) {
	vm.MEM[vm.SP] = data.Integer(0)
	vm.SP--
	vm.PC++
}
