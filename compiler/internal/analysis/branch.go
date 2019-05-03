package analysis

import "gitlab.com/kode4food/ale/runtime/isa"

type (
	// Node is returned when a Branch analysis is performed
	Node interface {
		Base() isa.Offset
		Code() isa.Instructions
	}

	// Instructions represents a series of non-branching instructions
	Instructions interface {
		Node
		Set(isa.Instructions)
	}

	// Branches represents a branching junction
	Branches interface {
		Node
		Prologue() Instructions
		ThenBranch() Node
		ElseBranch() Node
		Epilogue() Node
	}

	instructions struct {
		base isa.Offset
		code isa.Instructions
	}

	branches struct {
		base       isa.Offset
		prologue   Instructions
		thenBranch Node
		elseBranch Node
		epilogue   Node
	}
)

// Branch performs conditional branch analysis
func Branch(code isa.Instructions) Node {
	return branchFrom(0, code)
}

func branchFrom(base isa.Offset, code isa.Instructions) Node {
	for pc, inst := range code {
		oc := inst.Opcode
		if oc != isa.CondJump {
			continue
		}
		if rs := splitCond(base+isa.Offset(pc), code[pc:]); rs != nil {
			rs.base = base
			rs.prologue = &instructions{
				base: base,
				code: code[0 : pc+1],
			}
			return rs
		}
	}
	return &instructions{
		base: base,
		code: code,
	}
}

func splitCond(base isa.Offset, code isa.Instructions) *branches {
	thenIdx := isa.Index(code[0].Args[0])
	thenLabel := findLabel(code, thenIdx)
	if thenLabel <= 0 {
		return nil // not part of this block
	}

	prev := code[thenLabel-1]
	if prev.Opcode != isa.Jump {
		return nil // not created with build.Cond
	}

	elseRes := code[1:thenLabel]
	endLabel := findLabel(code, isa.Index(prev.Args[0]))
	thenRes := code[thenLabel:endLabel]
	return &branches{
		thenBranch: branchFrom(base+isa.Offset(thenLabel), thenRes),
		elseBranch: branchFrom(base+1, elseRes),
		epilogue:   branchFrom(base+isa.Offset(endLabel), code[endLabel:]),
	}
}

func (i *instructions) Base() isa.Offset {
	return i.base
}

func (i *instructions) Set(code isa.Instructions) {
	i.code = code
}

func (i *instructions) Code() isa.Instructions {
	return i.code
}

func (b *branches) Base() isa.Offset {
	return b.base
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
	res = append(res, b.thenBranch.Code()...)
	res = append(res, b.epilogue.Code()...)
	return res
}
