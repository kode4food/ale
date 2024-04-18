package vm

import (
	"github.com/kode4food/ale/pkg/data"
	"github.com/kode4food/ale/pkg/env"
)

func doArg(vm *machine) {
	vm.MEM[vm.SP] = vm.ARGS[vm.INST.Operand()]
	vm.SP--
	vm.PC++
}

func doArgLen(vm *machine) {
	vm.MEM[vm.SP] = data.Integer(len(vm.ARGS))
	vm.SP--
	vm.PC++
}

func doBind(vm *machine) {
	vm.SP++
	name := vm.MEM[vm.SP].(data.Local)
	vm.SP++
	vm.CL.Globals.Declare(name).Bind(vm.MEM[vm.SP])
	vm.PC++
}

func doBindRef(vm *machine) {
	vm.SP++
	ref := vm.MEM[vm.SP].(*Ref)
	vm.SP++
	ref.Value = vm.MEM[vm.SP]
	vm.PC++
}

func doClosure(vm *machine) {
	vm.MEM[vm.SP] = vm.CL.captured[vm.INST.Operand()]
	vm.SP--
	vm.PC++
}

func doDeclare(vm *machine) {
	vm.SP++
	vm.CL.Globals.Declare(
		vm.MEM[vm.SP].(data.Local),
	)
	vm.PC++
}

func doDeref(vm *machine) {
	SP1 := vm.SP + 1
	vm.MEM[SP1] = vm.MEM[SP1].(*Ref).Value
	vm.PC++
}

func doDup(vm *machine) {
	vm.MEM[vm.SP] = vm.MEM[vm.SP+1]
	vm.SP--
	vm.PC++
}

func doLoad(vm *machine) {
	vm.MEM[vm.SP] = vm.MEM[vm.LP+int(vm.INST.Operand())]
	vm.SP--
	vm.PC++
}

func doNewRef(vm *machine) {
	vm.MEM[vm.SP] = new(Ref)
	vm.SP--
	vm.PC++
}

func doPop(vm *machine) {
	vm.SP++
	vm.PC++
}

func doPrivate(vm *machine) {
	vm.SP++
	vm.CL.Globals.Private(
		vm.MEM[vm.SP].(data.Local),
	)
	vm.PC++
}

func doResolve(vm *machine) {
	SP1 := vm.SP + 1
	vm.MEM[SP1] = env.MustResolveValue(
		vm.CL.Globals,
		vm.MEM[SP1].(data.Symbol),
	)
	vm.PC++
}

func doRestArg(vm *machine) {
	vm.MEM[vm.SP] = vm.ARGS[vm.INST.Operand():]
	vm.SP--
	vm.PC++
}

func doSetArgs(vm *machine) {
	vm.SP++
	vm.ARGS = vm.MEM[vm.SP].(data.Vector)
	vm.PC++
}

func doStore(vm *machine) {
	vm.SP++
	vm.MEM[vm.LP+int(vm.INST.Operand())] = vm.MEM[vm.SP]
	vm.PC++
}
