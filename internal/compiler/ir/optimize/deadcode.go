package optimize

import (
	"github.com/kode4food/ale/internal/compiler/encoder"
	"github.com/kode4food/ale/internal/compiler/ir/visitor"
	"github.com/kode4food/ale/internal/runtime/isa"
	"github.com/kode4food/comb/basics"
)

var ineffectivePushes = globalRepeatedReplace(
	visitor.Pattern{
		selectEffects(func(e *isa.Effect) bool {
			return e.Push == 1 && e.Pop == 0 && !e.DPop && !e.Exit
		}),
		{isa.Pop},
	},
	removeInstructions,
)

func ineffectiveStores(e *encoder.Encoded) {
	for {
		if res := replaceIneffectiveStore(e.Code); res != nil {
			e.Code = res
			continue
		}
		break
	}
}

func replaceIneffectiveStore(c isa.Instructions) isa.Instructions {
	idx := findOpcode(c, isa.Store)
	if idx == -1 {
		return nil
	}
	store := c[idx]
	op := store.Operand()
	rest := c[idx+1:]
	if len(rest) == 0 {
		return nil
	}
	if rest[0].Opcode() != isa.Load || rest[0].Operand() != op {
		if res := replaceIneffectiveStore(rest); res != nil {
			return append(c[0:idx+1], res...)
		}
		return nil
	}
	if len(basics.Filter(rest[1:], func(i isa.Instruction) bool {
		oc := i.Opcode()
		return (oc == isa.Store || oc == isa.Load) && i.Operand() == op
	})) > 0 {
		if res := replaceIneffectiveStore(rest); res != nil {
			return append(c[0:idx+1], res...)
		}
		return nil
	}
	return append(c[0:idx], rest[1:]...)
}

func findOpcode(c isa.Instructions, oc isa.Opcode) int {
	for idx, i := range c {
		if i.Opcode() == oc {
			return idx
		}
	}
	return -1
}
