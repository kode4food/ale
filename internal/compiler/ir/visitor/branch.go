package visitor

import (
	"github.com/kode4food/ale/internal/basics"
	"github.com/kode4food/ale/internal/runtime/isa"
)

type (
	// A Node is returned when a Branched analysis is performed
	Node interface {
		// Code returns the instructions for this node
		Code() isa.Instructions

		// isModified returns whether this node has been modified
		isModified() bool
	}

	// Instructions represent a series of non-branching instructions
	Instructions interface {
		Node

		// Set updates the instructions for this node
		Set(isa.Instructions)
	}

	// Branches represent a branching junction
	Branches interface {
		Node

		// Prologue returns the instructions before the branching
		Prologue() Instructions

		// ThenBranch returns the then branch node
		ThenBranch() Node

		// ElseBranch returns the else branch node
		ElseBranch() Node

		// Epilogue returns the instructions after the branching
		Epilogue() Node
	}

	instructions struct {
		code  isa.Instructions
		dirty bool
	}

	BranchScanner struct {
		Then     Scanner
		Else     Scanner
		Epilogue Scanner
	}

	Scanner func(isa.Instructions) Node

	branches struct {
		prologue   Instructions
		elseBranch Node
		thenBranch Node
		epilogue   Node
		code       isa.Instructions
		elseJump   isa.Instruction
		thenLabel  isa.Instruction
		joinLabel  isa.Instruction
	}
)

// All treats a set of instructions as a single block (no branching)
func All(code isa.Instructions) Node {
	return &instructions{
		code: code,
	}
}

// Branched splits linear instructions into a tree of conditional branches
func Branched(code isa.Instructions) Node {
	b := &BranchScanner{
		Then:     Branched,
		Else:     Branched,
		Epilogue: Branched,
	}
	return b.Scan(code)
}

func (b *BranchScanner) Scan(code isa.Instructions) Node {
	for pc, inst := range code {
		if oc := inst.Opcode(); oc != isa.CondJump {
			continue
		}
		if rs := b.splitCondJump(code, pc); rs != nil {
			return rs
		}
	}
	return &instructions{
		code: code,
	}
}

func (b *BranchScanner) splitCondJump(
	code isa.Instructions, condJumpIdx int,
) *branches {
	prologue := &instructions{
		code: code[0 : condJumpIdx+1],
	}

	condJump := code[condJumpIdx]
	rest := code[condJumpIdx+1:]
	idx := condJump.Operand()
	thenIdx, thenLabel := findLabel(rest, idx)
	if thenIdx <= 0 {
		return nil // not part of this block
	}

	elseJumpIdx := thenIdx - 1
	elseJump := rest[elseJumpIdx]
	oc, idx := elseJump.Split()
	if oc != isa.Jump { // jump expected for generated branches
		return nil
	}

	joinIdx, joinLabel := findLabel(rest, idx)
	if joinIdx <= thenIdx { // forward jump expected in generated branches
		return nil
	}

	return &branches{
		code:       code,
		prologue:   prologue,
		elseBranch: b.Else(rest[:elseJumpIdx]),
		elseJump:   elseJump,
		thenLabel:  thenLabel,
		thenBranch: b.Then(rest[thenIdx+1 : joinIdx]),
		joinLabel:  joinLabel,
		epilogue:   b.Epilogue(rest[joinIdx+1:]),
	}
}

func (i *instructions) Set(code isa.Instructions) {
	if !i.dirty && !basics.Equal(i.code, code) {
		i.dirty = true
	}
	i.code = code
}

func (i *instructions) Code() isa.Instructions {
	return i.code
}

func (i *instructions) isModified() bool {
	return i.dirty
}

func (b *branches) Prologue() Instructions {
	return b.prologue
}

func (b *branches) ThenBranch() Node {
	return b.thenBranch
}

func (b *branches) ElseBranch() Node {
	return b.elseBranch
}

func (b *branches) Epilogue() Node {
	return b.epilogue
}

func (b *branches) Code() isa.Instructions {
	if !b.isModified() {
		return b.code
	}
	res := isa.Instructions{}
	res = append(res, b.prologue.Code()...)
	res = append(res, b.elseBranch.Code()...)
	res = append(res, b.elseJump)
	res = append(res, b.thenLabel)
	res = append(res, b.thenBranch.Code()...)
	res = append(res, b.joinLabel)
	res = append(res, b.epilogue.Code()...)
	return res
}

func (b *branches) isModified() bool {
	return b.prologue.isModified() ||
		b.elseBranch.isModified() ||
		b.thenBranch.isModified() ||
		b.epilogue.isModified()
}

func findLabel(code isa.Instructions, lbl isa.Operand) (int, isa.Instruction) {
	for pc, inst := range code {
		if oc, op := inst.Split(); oc == isa.Label && op == lbl {
			return pc, inst
		}
	}
	return -1, 0
}
