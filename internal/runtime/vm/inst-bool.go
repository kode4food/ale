package vm

import "github.com/kode4food/ale/pkg/data"

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
