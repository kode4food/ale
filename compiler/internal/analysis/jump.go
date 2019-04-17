package analysis

import (
	"fmt"

	"gitlab.com/kode4food/ale/runtime/isa"
)

func verifyJumps(code isa.Instructions) {
	for _, l := range code {
		oc := l.Opcode
		if oc == isa.CondJump || oc == isa.Jump {
			mustFindLabel(code, isa.Index(l.Args[0]))
		}
	}
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
