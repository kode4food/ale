package analysis

import (
	"errors"
	"fmt"

	"github.com/kode4food/ale/internal/compiler/ir/visitor"
	"github.com/kode4food/ale/internal/runtime/isa"
)

type stackSizes struct {
	maxSize int
	endSize int
}

const (
	// ErrBadStackTermination is raised when the analyzer verifies a branch
	// that ends in a non-empty state
	ErrBadStackTermination = "invalid stack end-state: %d"

	// ErrBadBranchTermination is raised when the analyzer verifies parallel
	// branches that end with a different stack size
	ErrBadBranchTermination = "branches should end in the same state"
)

func verifyStackSize(code isa.Instructions) error {
	s := new(stackSizes)
	if err := s.calculateNode(visitor.Branched(code)); err != nil {
		return err
	}
	if s.endSize != 0 {
		return fmt.Errorf(ErrBadStackTermination, s.endSize)
	}
	return nil
}

// CalculateStackSize returns the maximum depth for the stack based on the
// Instructions provided.
func CalculateStackSize(code isa.Instructions) (isa.Operand, error) {
	s := new(stackSizes)
	if err := s.calculateNode(visitor.Branched(code)); err != nil {
		return 0, err
	}
	return isa.Operand(s.maxSize), nil
}

func (s *stackSizes) calculateNode(n visitor.Node) error {
	switch n := n.(type) {
	case visitor.Branches:
		if err := s.calculateInstructions(n.Prologue()); err != nil {
			return err
		}
		t := n.ThenBranch()
		e := n.ElseBranch()
		if err := s.calculateBranches(t, e); err != nil {
			return err
		}
		return s.calculateNode(n.Epilogue())
	case visitor.Instructions:
		return s.calculateInstructions(n)
	default:
		return nil
	}
}

func (s *stackSizes) calculateInstructions(inst visitor.Instructions) error {
	for _, inst := range inst.Code() {
		if err := s.calculateInstruction(inst); err != nil {
			return err
		}
	}
	return nil
}

func (s *stackSizes) calculateInstruction(inst isa.Instruction) error {
	c, err := inst.StackChange()
	if err != nil {
		return err
	}
	s.endSize += c
	s.maxSize = max(s.endSize, s.maxSize)
	return nil
}

func (s *stackSizes) calculateBranches(thenNode, elseNode visitor.Node) error {
	thenRes, err := s.calculateBranch(thenNode)
	if err != nil {
		return err
	}
	elseRes, err := s.calculateBranch(elseNode)
	if err != nil {
		return err
	}
	if elseRes.endSize != thenRes.endSize {
		return errors.New(ErrBadBranchTermination)
	}
	s.endSize += elseRes.endSize
	return nil
}

func (s *stackSizes) calculateBranch(n visitor.Node) (*stackSizes, error) {
	res := &stackSizes{
		maxSize: s.maxSize,
	}
	if err := res.calculateNode(n); err != nil {
		return nil, err
	}
	s.maxSize = max(s.maxSize, res.maxSize)
	return res, nil
}
