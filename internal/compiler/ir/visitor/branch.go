package visitor

import "github.com/kode4food/ale/internal/runtime/isa"

type (
	// A Node is returned when a Branched analysis is performed
	Node interface {
		Code() isa.Instructions
	}

	// Instructions represent a series of non-branching instructions
	Instructions interface {
		Node
		Set(isa.Instructions)
	}

	// Branches represent a branching junction
	Branches interface {
		Node
		Prologue() Instructions
		ThenBranch() Node
		ElseBranch() Node
		Epilogue() Node
	}

	instructions struct {
		code isa.Instructions
	}

	branches struct {
		prologue   Instructions
		elseBranch Node
		thenBranch Node
		epilogue   Node
		elseJump   isa.Instruction
		thenLabel  isa.Instruction
		joinLabel  isa.Instruction
	}
)

// All treats a set of instructions as a single block (no branching)
func All(code isa.Instructions) Instructions {
	return &instructions{
		code: code,
	}
}

// Branched splits linear instructions into a tree of conditional branches
func Branched(code isa.Instructions) Node {
	for pc, inst := range code {
		if oc := inst.Opcode(); oc != isa.CondJump {
			continue
		}
		if rs := splitCondJump(code, pc); rs != nil {
			return rs
		}
	}
	return &instructions{
		code: code,
	}
}

func splitCondJump(code isa.Instructions, condJumpIdx int) *branches {
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
		prologue:   prologue,
		elseBranch: Branched(rest[0:elseJumpIdx]),
		elseJump:   elseJump,
		thenLabel:  thenLabel,
		thenBranch: Branched(rest[thenIdx+1 : joinIdx]),
		joinLabel:  joinLabel,
		epilogue:   Branched(rest[joinIdx+1:]),
	}
}

func (i *instructions) Set(code isa.Instructions) {
	i.code = code
}

func (i *instructions) Code() isa.Instructions {
	return i.code
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

func findLabel(code isa.Instructions, lbl isa.Operand) (int, isa.Instruction) {
	for pc, inst := range code {
		if oc, op := inst.Split(); oc == isa.Label && op == lbl {
			return pc, inst
		}
	}
	return -1, 0
}
