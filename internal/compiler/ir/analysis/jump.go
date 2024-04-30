package analysis

import (
	"fmt"

	"github.com/kode4food/ale/internal/runtime/isa"
)

const (
	// ErrLabelNotAnchored is raised when a label hasn't been placed in an
	// isa.Instructions stream
	ErrLabelNotAnchored = "label not anchored: %d"

	// ErrLabelMultipleAnchors is raised when a label has been placed more than
	// once in an isa.Instructions stream
	ErrLabelMultipleAnchors = "label anchored multiple times: %d"
)

func verifyJumps(code isa.Instructions) error {
	for _, l := range code {
		if oc, op := l.Split(); oc == isa.CondJump || oc == isa.Jump {
			if _, err := findLabel(code, op); err != nil {
				return err
			}
		}
	}
	return nil
}

func findLabel(code isa.Instructions, lbl isa.Operand) (int, error) {
	res := -1
	for pc, inst := range code {
		if oc, op := inst.Split(); oc == isa.Label && op == lbl {
			if res != -1 {
				return res, fmt.Errorf(ErrLabelMultipleAnchors, lbl)
			}
			res = pc
		}
	}
	if res == -1 {
		return res, fmt.Errorf(ErrLabelNotAnchored, lbl)
	}
	return res, nil
}
