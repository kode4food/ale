package analysis

import (
	"fmt"

	"github.com/kode4food/ale/runtime/isa"
)

// Error messages
const (
	errLabelNotAnchored     = "label not anchored: %d"
	errLabelMultipleAnchors = "label anchored multiple times: %d"
)

func verifyJumps(code isa.Instructions) {
	for _, l := range code {
		oc := l.Opcode
		if oc == isa.CondJump || oc == isa.Jump {
			mustFindLabel(code, isa.Index(l.Args[0]))
		}
	}
}

func findLabel(code isa.Instructions, lbl isa.Index) (int, error) {
	ic := lbl.Word()
	res := -1
	for pc, inst := range code {
		if inst.Opcode == isa.Label && inst.Args[0] == ic {
			if res != -1 {
				return res, fmt.Errorf(errLabelMultipleAnchors, lbl)
			}
			res = pc
		}
	}
	if res == -1 {
		return res, fmt.Errorf(errLabelNotAnchored, lbl)
	}
	return res, nil
}

func mustFindLabel(code isa.Instructions, lbl isa.Index) int {
	res, err := findLabel(code, lbl)
	if err != nil {
		panic(err)
	}
	return res
}
