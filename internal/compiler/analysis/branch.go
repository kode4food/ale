package analysis

import (
	"fmt"

	"gitlab.com/kode4food/ale/internal/runtime/isa"
)

type (
	// Branches is returned when a Branch analysis is performed
	Branches struct {
		Prologue   isa.Instructions
		ThenBranch isa.Instructions
		ElseBranch isa.Instructions
		Epilogue   isa.Instructions
	}

	splitResult struct {
		thenBranch isa.Instructions
		elseBranch isa.Instructions
		epilogue   isa.Instructions
	}
)

// Branch performs conditional branch analysis
func Branch(code isa.Instructions) *Branches {
	for pc, inst := range code {
		oc := inst.Opcode
		if oc != isa.CondJump {
			continue
		}
		if rs, ok := splitCondJump(code[pc:]); ok {
			return &Branches{
				Prologue:   code[0 : pc+1],
				ThenBranch: rs.thenBranch,
				ElseBranch: rs.elseBranch,
				Epilogue:   rs.epilogue,
			}
		}
	}
	return &Branches{Prologue: code}
}

func findLabel(code isa.Instructions, lbl isa.Index) int {
	ic := isa.Word(lbl)
	for pc, inst := range code {
		if inst.Opcode == isa.Label && inst.Args[0] == ic {
			return pc
		}
	}
	return -1
}

func mustFindLabel(code isa.Instructions, lbl isa.Index) int {
	if res := findLabel(code, lbl); res >= 0 {
		return res
	}
	panic(fmt.Sprintf("label not anchored: %d", lbl))
}

func splitCondJump(code isa.Instructions) (*splitResult, bool) {
	thenIdx := isa.Index(code[0].Args[0])
	thenLabel := findLabel(code, thenIdx)
	if thenLabel <= 0 {
		return nil, false // not part of this block
	}

	prev := code[thenLabel-1]
	if prev.Opcode != isa.Jump {
		return nil, false // not created with build.Cond
	}

	elseRes := code[1:thenLabel]
	endLabel := findLabel(code, isa.Index(prev.Args[0]))
	thenRes := code[thenLabel:endLabel]
	return &splitResult{
		thenBranch: thenRes,
		elseBranch: elseRes,
		epilogue:   code[endLabel:],
	}, true
}
