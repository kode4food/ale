package vm

import (
	"slices"

	"github.com/kode4food/ale/pkg/data"
)

func doCar(vm *machine) {
	SP1 := vm.SP + 1
	vm.MEM[SP1] = vm.MEM[SP1].(data.Pair).Car()
	vm.PC++
}

func doCdr(vm *machine) {
	SP1 := vm.SP + 1
	vm.MEM[SP1] = vm.MEM[SP1].(data.Pair).Cdr()
	vm.PC++
}

func doCons(vm *machine) {
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

func doEmpty(vm *machine) {
	SP1 := vm.SP + 1
	vm.MEM[SP1] = data.Bool(vm.MEM[SP1].(data.Sequence).IsEmpty())
	vm.PC++
}

func doVector(vm *machine) {
	op := vm.INST.Operand()
	RES := vm.SP + int(op)
	vm.MEM[RES] = slices.Clone(vm.MEM[vm.SP+1 : RES+1])
	vm.SP = RES - 1
	vm.PC++
}
