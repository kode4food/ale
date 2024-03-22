package analysis

import (
	"errors"
	"fmt"

	"github.com/kode4food/ale/compiler/ir/visitor"
	"github.com/kode4food/ale/runtime/isa"
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

// MustCalculateStackSize is a wrapper around CalculateStackSize that will
// panic if the Instructions provided are invalid
func MustCalculateStackSize(code isa.Instructions) isa.Operand {
	res, err := CalculateStackSize(code)
	if err != nil {
		panic(err)
	}
	return res
}

func (s *stackSizes) calculateNode(n visitor.Node) error {
	switch n := n.(type) {
	case visitor.Branches:
		s.calculateInstructions(n.Prologue())
		t := n.ThenBranch()
		e := n.ElseBranch()
		if err := s.calculateBranches(t, e); err != nil {
			return err
		}
		if err := s.calculateNode(n.Epilogue()); err != nil {
			return err
		}
	case visitor.Instructions:
		s.calculateInstructions(n)
	}
	return nil
}

func (s *stackSizes) calculateInstructions(inst visitor.Instructions) {
	for _, inst := range inst.Code() {
		s.calculateInstruction(inst)
	}
}

func (s *stackSizes) calculateInstruction(inst isa.Instruction) {
	s.endSize += inst.StackChange()
	s.maxSize = max(s.endSize, s.maxSize)
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
