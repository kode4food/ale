package visitor

import (
	"bytes"
	"strings"

	"gitlab.com/kode4food/ale/runtime/isa"
)

type (
	// Node is returned when a Branch analysis is performed
	Node interface {
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
		code isa.Instructions
	}

	branches struct {
		prologue   Instructions
		elseBranch Node
		elseJump   *isa.Instruction
		thenLabel  *isa.Instruction
		thenBranch Node
		joinLabel  *isa.Instruction
		epilogue   Node
	}
)

// Branch splits linear instructions into a tree conditional branches
func Branch(code isa.Instructions) Node {
	for pc, inst := range code {
		oc := inst.Opcode
		if oc != isa.CondJump {
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
	thenIdx, thenLabel := findLabel(rest, isa.Index(condJump.Args[0]))
	if thenIdx <= 0 {
		return nil // not part of this block
	}

	elseJumpIdx := thenIdx - 1
	elseJump := rest[elseJumpIdx]
	if elseJump.Opcode != isa.Jump {
		return nil // not created with build.Cond
	}

	joinIdx, joinLabel := findLabel(rest, isa.Index(elseJump.Args[0]))

	return &branches{
		prologue:   prologue,
		elseBranch: Branch(rest[0:elseJumpIdx]),
		elseJump:   elseJump,
		thenLabel:  thenLabel,
		thenBranch: Branch(rest[thenIdx+1 : joinIdx]),
		joinLabel:  joinLabel,
		epilogue:   Branch(rest[joinIdx+1:]),
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

func (b *branches) String() string {
	return indentedString(0, b)
}

func (i *instructions) String() string {
	return indentedString(0, i)
}

func findLabel(code isa.Instructions, lbl isa.Index) (int, *isa.Instruction) {
	ic := lbl.Word()
	for pc, inst := range code {
		if inst.Opcode == isa.Label && inst.Args[0] == ic {
			return pc, inst
		}
	}
	return -1, nil
}

func indentedString(lvl int, n Node) string {
	var buf bytes.Buffer
	switch typed := n.(type) {
	case Branches:
		buf.WriteString(indentedString(lvl, typed.Prologue()))
		buf.WriteString(indentedString(lvl+1, typed.ThenBranch()))
		buf.WriteString(strings.Repeat("  ", lvl))
		buf.WriteString("else:\n")
		buf.WriteString(indentedString(lvl+1, typed.ElseBranch()))
		buf.WriteString(indentedString(lvl, typed.Epilogue()))
	case Instructions:
		for _, i := range typed.Code() {
			buf.WriteString(strings.Repeat("  ", lvl))
			buf.WriteString(i.String())
			buf.WriteString("\n")
		}
	}
	return buf.String()
}
