package vm

import "github.com/kode4food/ale/data"

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
