package optimize

import (
	"github.com/kode4food/ale/internal/compiler/encoder"
	"github.com/kode4food/ale/internal/compiler/ir/visitor"
	"github.com/kode4food/ale/internal/runtime/isa"
	"github.com/kode4food/comb/basics"
)

// ineffectivePushes deletes values pushed to the stack for no reason,
// specifically if they're literals followed immediately by a pop instruction
var ineffectivePushes = globalReplace(
	visitor.Pattern{
		selectEffects(func(e *isa.Effect) bool {
			return e.Push == 1 && e.Pop == 0 && !e.DPop && !e.Exit
		}),
		{isa.Pop},
	},
	removeInstructions,
)

// ineffectiveStores deletes isolated store instructions followed by a load
// instruction wit the same operand
func ineffectiveStores(e *encoder.Encoded) {
	if res := replaceIneffectiveStore(e.Code); res != nil {
		e.Code = res
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

	if rest[0] != isa.Load.New(op) || hasConflictingLoadStore(rest[1:], op) {
		if res := replaceIneffectiveStore(rest); res != nil {
			return append(c[0:idx+1], res...)
		}
		return nil
	}
	return append(c[0:idx], rest[1:]...)
}

func hasConflictingLoadStore(c isa.Instructions, op isa.Operand) bool {
	load := isa.Load.New(op)
	store := isa.Store.New(op)
	return len(basics.Filter(c, func(i isa.Instruction) bool {
		return i == load || i == store
	})) > 0
}

func findOpcode(c isa.Instructions, oc isa.Opcode) int {
	for idx, i := range c {
		if i.Opcode() == oc {
			return idx
		}
	}
	return -1
}
