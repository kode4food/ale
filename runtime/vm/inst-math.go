package vm

import "github.com/kode4food/ale/data"

func doAdd(vm *VM) {
	vm.SP++
	SP1 := vm.SP + 1
	vm.MEM[SP1] = vm.MEM[SP1].(data.Number).Add(
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

func doMul(vm *VM) {
	vm.SP++
	SP1 := vm.SP + 1
	vm.MEM[SP1] = vm.MEM[SP1].(data.Number).Mul(
		vm.MEM[vm.SP].(data.Number),
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

func doSub(vm *VM) {
	vm.SP++
	SP1 := vm.SP + 1
	vm.MEM[SP1] = vm.MEM[SP1].(data.Number).Sub(
		vm.MEM[vm.SP].(data.Number),
	)
	vm.PC++
}
