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

func verifyStackSize(code isa.Instructions) {
	s := new(stackSizes)
	s.calculateNode(visitor.Branch(code))
	if s.endSize != 0 {
		panic(fmt.Errorf(ErrBadStackTermination, s.endSize))
	}
}

// CalculateStackSize returns the maximum and final depths for the stack based
// on the instructions provided. If the final depth is non-zero, this is
// usually an indication that bad instructions were encoded
func CalculateStackSize(code isa.Instructions) (isa.Operand, isa.Operand) {
	s := new(stackSizes)
	s.calculateNode(visitor.Branch(code))
	return isa.Operand(s.maxSize), isa.Operand(s.endSize)
}

func (s *stackSizes) calculateNode(n visitor.Node) {
	switch n := n.(type) {
	case visitor.Branches:
		s.calculateInstructions(n.Prologue())
		s.calculateBranches(n.ThenBranch(), n.ElseBranch())
		s.calculateNode(n.Epilogue())
	case visitor.Instructions:
		s.calculateInstructions(n)
	}
}

func (s *stackSizes) calculateInstructions(inst visitor.Instructions) {
	for _, inst := range inst.Code() {
		s.calculateInstruction(inst)
	}
}

func (s *stackSizes) calculateInstruction(inst isa.Instruction) {
	effect := isa.MustGetEffect(inst.Opcode())
	dPop := getStackChange(inst, effect.DPop)
	s.endSize += (effect.Push - effect.Pop) - dPop
	s.maxSize = max(s.endSize, s.maxSize)
}

func (s *stackSizes) calculateBranches(thenNode, elseNode visitor.Node) {
	thenRes := s.calculateBranch(thenNode)
	elseRes := s.calculateBranch(elseNode)
	if elseRes.endSize != thenRes.endSize {
		panic(errors.New(ErrBadBranchTermination))
	}
	s.endSize += elseRes.endSize
}

func (s *stackSizes) calculateBranch(n visitor.Node) *stackSizes {
	res := &stackSizes{
		maxSize: s.maxSize,
	}
	res.calculateNode(n)
	s.maxSize = max(s.maxSize, res.maxSize)
	return res
}

func getStackChange(inst isa.Instruction, dPop bool) int {
	if dPop {
		return int(inst.Operand())
	}
	return 0
}
