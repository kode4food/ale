package optimize

import (
	"slices"

	"github.com/kode4food/ale/internal/basics"
	"github.com/kode4food/ale/internal/runtime/isa"
)

func hasAnyArgInstruction(c isa.Instructions) bool {
	return slices.ContainsFunc(c, isArgInstruction)
}

func filterArgInstructions(c isa.Instructions) isa.Instructions {
	return basics.Filter(c, isArgInstruction)
}

func isArgInstruction(i isa.Instruction) bool {
	switch i.Opcode() {
	case isa.ArgsPush, isa.ArgsPop, isa.Arg, isa.ArgsLen, isa.ArgsRest:
		return true
	default:
		return false
	}
}

func canMapArgsToLocals(c isa.Instructions, argc isa.Operand) bool {
	highArg := -1
	for _, i := range filterArgInstructions(c) {
		switch i.Opcode() {
		case isa.Arg:
			idx := int(i.Operand())
			if idx > highArg {
				highArg = idx
			}
		case isa.ArgsPush, isa.ArgsPop, isa.ArgsLen, isa.ArgsRest:
			return false
		default:
			// no-op
		}
	}
	return highArg < int(argc)
}

func mapArgsToLocals(
	c isa.Instructions, argsBase, argc isa.Operand,
) isa.Instructions {
	al := makeArgLocalMap(c, argsBase)
	res := make(isa.Instructions, 0, len(c)+int(argc))
	for i := range int(argc) {
		if to, ok := al[isa.Operand(i)]; ok {
			res = append(res, isa.Store.New(to))
			continue
		}
		res = append(res, isa.Pop.New())
	}
	res = append(res, basics.Map(c, func(i isa.Instruction) isa.Instruction {
		if i.Opcode() == isa.Arg {
			to := al[i.Operand()]
			return isa.Load.New(to)
		}
		return i
	})...)
	return res
}

func makeArgLocalMap(c isa.Instructions, argsBase isa.Operand) operandMap {
	next := argsBase
	res := operandMap{}
	for _, i := range c {
		if i.Opcode() != isa.Arg {
			continue
		}
		op := i.Operand()
		if _, ok := res[op]; !ok {
			res[op] = next
			next++
		}
	}
	return res
}

func stackArgs(c isa.Instructions, argc isa.Operand) isa.Instructions {
	res := make(isa.Instructions, 0, len(c)+2)
	res = append(res, isa.ArgsPush.New(argc))
	res = append(res, c...)
	res = append(res, isa.ArgsPop.New())
	return res
}
