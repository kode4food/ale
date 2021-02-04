package analysis

import (
	"fmt"

	"github.com/kode4food/ale/runtime/isa"
)

// Error messages
const (
	ErrLabelNotAnchored = "label not anchored: %d"
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
	ic := lbl.Word()
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
	panic(fmt.Errorf(ErrLabelNotAnchored, lbl))
}
