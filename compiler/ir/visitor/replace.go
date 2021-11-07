package visitor

import "github.com/kode4food/ale/runtime/isa"

type (
	// Pattern is a replacement pattern for the visitor. The
	// second level of the array identifies the possible matching
	// Opcodes, while the first level identifies the sequence in
	// which they should appear
	Pattern [][]isa.Opcode

	// Mapper maps one set of instructions to another
	Mapper func(isa.Instructions) isa.Instructions

	replace struct {
		pattern Pattern
		mapper  Mapper
	}
)

// Replace visits all Instruction nodes and if any of the instructions therein
// match the provided Pattern, they will be replaced using the provided Mapper
func Replace(root Node, pattern Pattern, mapper Mapper) {
	r := &replace{
		pattern: pattern,
		mapper:  mapper,
	}
	DepthFirst(root, r)
}

func (*replace) EnterRoot(Node)         {}
func (*replace) ExitRoot(Node)          {}
func (*replace) EnterBranches(Branches) {}
func (*replace) ExitBranches(Branches)  {}

func (r *replace) Instructions(i Instructions) {
	pattern := r.pattern
	code := i.Code()
	var state, start, found int
	res := isa.Instructions{}
	for pc, inst := range code {
		if pattern.matchesState(inst.Opcode, state) {
			if state == 0 {
				found = pc
			}
			state++
			if state < len(pattern) {
				continue
			}
			next := pc + 1
			res = append(res, code[start:found]...)
			res = append(res, r.mapper(code[found:next])...)
			start = next
		}
		state = 0
	}
	res = append(res, code[start:]...)
	i.Set(res)
}

func (p Pattern) matchesState(opcode isa.Opcode, state int) bool {
	set := p[state]
	for _, elem := range set {
		if elem == opcode {
			return true
		}
	}
	return false
}
