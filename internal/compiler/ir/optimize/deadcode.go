package optimize

import (
	"slices"

	"github.com/kode4food/ale/internal/compiler/encoder"
	"github.com/kode4food/ale/internal/compiler/ir/visitor"
	"github.com/kode4food/ale/internal/runtime/isa"
	"github.com/kode4food/comb/basics"
)

var (
	deadCode = repeatWhenModified(
		ineffectiveStores, ineffectiveCopies, ineffectivePushes,
	)

	// ineffectivePushes deletes values pushed to the stack for no reason,
	// specifically if they're literals followed immediately by a pop instruction
	ineffectivePushes = globalReplace(
		visitor.Pattern{
			selectEffects(func(e *isa.Effect) bool {
				return e.Push == 1 && e.Pop == 0 && !e.DPop && !e.Exit
			}),
			{isa.Pop},
		},
		removeInstructions,
	)
)

// ineffectiveStores deletes isolated store instructions followed by a load
// instruction with the same operand
func ineffectiveStores(e *encoder.Encoded) *encoder.Encoded {
	if c := replaceIneffectiveStore(e.Code); c != nil {
		e.Code = c
	}
	return e
}

func replaceIneffectiveStore(c isa.Instructions) isa.Instructions {
	idx := findOpcode(c, isa.Store)
	if idx == -1 || hasReverseJump(c, idx) {
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

func ineffectiveCopies(e *encoder.Encoded) *encoder.Encoded {
	if c := replaceIneffectiveCopy(e.Code); c != nil {
		e.Code = c
	}
	return e
}

func replaceIneffectiveCopy(c isa.Instructions) isa.Instructions {
	idx := findOpcode(c, isa.Load)
	if idx == -1 || hasReverseJump(c, idx) {
		return nil
	}

	load := c[idx]
	rest := c[idx+1:]
	if len(rest) == 0 || rest[0].Opcode() != isa.Store {
		return nil
	}

	to := rest[0].Operand()
	if hasConflictingStore(rest[1:], load.Operand(), to) {
		return nil
	}

	return append(
		c[0:idx], mapIneffectiveLoads(rest[1:], isa.Load.New(to), load)...,
	)
}

func hasConflictingStore(c isa.Instructions, op ...isa.Operand) bool {
	return len(basics.Filter(c, func(i isa.Instruction) bool {
		return i.Opcode() == isa.Store && slices.Contains(op, i.Operand())
	})) > 0
}

func mapIneffectiveLoads(
	c isa.Instructions, from isa.Instruction, to isa.Instruction,
) isa.Instructions {
	return basics.Map(c, func(i isa.Instruction) isa.Instruction {
		if i == from {
			return to
		}
		return i
	})
}

func hasReverseJump(c isa.Instructions, before int) bool {
	offsets := map[isa.Operand]int{}
	for i, inst := range c {
		switch inst.Opcode() {
		case isa.Jump, isa.CondJump:
			if i <= before {
				continue
			}
			if o, ok := offsets[inst.Operand()]; ok && o < before {
				return true
			}

		case isa.Label:
			op := inst.Operand()
			if _, ok := offsets[op]; !ok {
				offsets[op] = i
			}

		default:
			// no-op
		}
	}
	return false
}
