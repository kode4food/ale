package vm

import (
	"errors"

	"github.com/kode4food/ale/pkg/data"
)

func doCondJump(vm *machine) {
	vm.SP++
	if vm.MEM[vm.SP] != data.False {
		vm.PC = int(vm.INST.Operand())
		return
	}
	vm.PC++
}

func doJump(vm *machine) {
	vm.PC = int(vm.INST.Operand())
}

func doNoOp(vm *machine) {
	vm.PC++
}

func doPanic(vm *machine) {
	vm.ST = failure
	panic(errors.New(data.ToString(vm.MEM[vm.SP+1])))
}

func doReturn(vm *machine) {
	vm.ST = success
	vm.RES = vm.MEM[vm.SP+1]
}

func doRetFalse(vm *machine) {
	vm.ST = success
	vm.RES = data.False
}

func doRetNull(vm *machine) {
	vm.ST = success
	vm.RES = data.Null
}

func doRetTrue(vm *machine) {
	vm.ST = success
	vm.RES = data.True
}
