package vm

import (
	"errors"

	"github.com/kode4food/ale/pkg/data"
)

func doCondJump(vm *VM) {
	vm.SP++
	if vm.MEM[vm.SP] != data.False {
		vm.PC = int(vm.INST.Operand())
		return
	}
	vm.PC++
}

func doJump(vm *VM) {
	vm.PC = int(vm.INST.Operand())
}

func doNoOp(vm *VM) {
	vm.PC++
}

func doPanic(vm *VM) {
	vm.ST = failure
	panic(errors.New(data.ToString(vm.MEM[vm.SP+1])))
}

func doReturn(vm *VM) {
	vm.ST = success
	vm.RES = vm.MEM[vm.SP+1]
}

func doRetFalse(vm *VM) {
	vm.ST = success
	vm.RES = data.False
}

func doRetNull(vm *VM) {
	vm.ST = success
	vm.RES = data.Null
}

func doRetTrue(vm *VM) {
	vm.ST = success
	vm.RES = data.True
}
