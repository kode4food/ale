package visitor

import (
	"slices"

	"github.com/kode4food/ale/pkg/runtime/isa"
)

type (
	// Pattern is a replacement pattern for the visitor. The second level of
	// the array identifies the possible matching Opcodes, while the first
	// level identifies the sequence in which they should appear
	Pattern [][]isa.Opcode

	// Mapper maps one set of instructions to another
	Mapper func(isa.Instructions) isa.Instructions

	Replacer struct {
		mapper  Mapper
		pattern Pattern
	}
)

const AnyOpcode = isa.OpcodeMask + 1

// Replace visits all Instruction nodes and if any of the instructions therein
// match the provided Pattern, they will be replaced using the provided Mapper
func Replace(pattern Pattern, mapper Mapper) *Replacer {
	return &Replacer{
		pattern: pattern,
		mapper:  mapper,
	}
}

func (*Replacer) EnterRoot(Node)         {}
func (*Replacer) ExitRoot(Node)          {}
func (*Replacer) EnterBranches(Branches) {}
func (*Replacer) ExitBranches(Branches)  {}

func (r *Replacer) Instructions(i Instructions) {
	pattern := r.pattern
	code := i.Code()
	var state, start, found int
	res := isa.Instructions{}
	for pc := 0; pc < len(code); {
		inst := code[pc]
		if oc := inst.Opcode(); !pattern.matchesState(oc, state) {
			if state == 0 {
				pc++
			} else {
				state = 0
			}
			continue
		}
		if state == 0 {
			found = pc
		}
		pc++
		if state < len(pattern)-1 {
			state++
			continue
		}
		res = append(res, code[start:found]...)
		res = append(res, r.mapper(code[found:pc])...)
		start = pc
		state = 0
	}
	res = append(res, code[start:]...)
	i.Set(res)
}

func (p Pattern) matchesState(opcode isa.Opcode, state int) bool {
	if len(p[state]) == 1 && p[state][0] == AnyOpcode {
		return true
	}
	return slices.Contains(p[state], opcode)
}
