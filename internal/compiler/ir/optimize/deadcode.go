package optimize

import (
	"slices"

	"github.com/kode4food/ale/internal/compiler/encoder"
	"github.com/kode4food/ale/internal/compiler/ir/visitor"
	"github.com/kode4food/ale/internal/runtime/isa"
	"github.com/kode4food/comb/basics"
)

var (
	deadCode = repeatWhenModified(redundantLocals, ineffectivePushes)

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

// redundantLocals deletes or rewrites Load and Store combinations that result
// in excessive memory location access or modification
func redundantLocals(e *encoder.Encoded) *encoder.Encoded {
	if c, ok := replaceRedundantLocals(e.Code); ok {
		e.Code = c
	}
	return e
}

func replaceRedundantLocals(c isa.Instructions) (isa.Instructions, bool) {
	labels := map[isa.Operand]int{}
	res := make(isa.Instructions, 0, len(c))
	var dirty bool

	for idx := 0; idx < len(c); {
		inst := c[idx]

		switch inst.Opcode() {
		case isa.Jump, isa.CondJump:
			if off, ok := labels[inst.Operand()]; ok && off < idx {
				return nil, false
			}

		case isa.Label:
			labels[inst.Operand()] = idx

		case isa.Store:
			if idx == len(c)-1 || c[idx+1].Opcode() != isa.Load {
				break
			}
			op := inst.Operand()
			next := c[idx+1]
			if next.Operand() != op || hasConflictingLoadStore(c[idx+2:], op) {
				break
			}
			c = slices.Concat(c[:idx], c[idx+2:])
			dirty = true
			continue

		case isa.Load:
			if idx == len(c)-1 || c[idx+1].Opcode() != isa.Store {
				break
			}
			next := c[idx+1]
			from := inst.Operand()
			to := next.Operand()
			if hasConflictingStore(c[idx+2:], from, to) {
				break
			}
			c = slices.Concat(
				c[:idx],
				mapIneffectiveLoads(c[idx+2:], isa.Load.New(to), inst),
			)
			dirty = true
			continue

		default:
			// no-op
		}
		res = append(res, inst)
		idx++
	}
	return res, dirty
}

func hasConflictingLoadStore(c isa.Instructions, op isa.Operand) bool {
	load := isa.Load.New(op)
	store := isa.Store.New(op)
	return len(basics.Filter(c, func(i isa.Instruction) bool {
		return i == load || i == store
	})) > 0
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
