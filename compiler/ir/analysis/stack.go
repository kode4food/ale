package analysis

import (
	"fmt"

	"github.com/kode4food/ale/compiler/ir/visitor"
	"github.com/kode4food/ale/internal/util"
	"github.com/kode4food/ale/runtime/isa"
)

type stackSizes struct {
	maxSize int
	endSize int
}

// Error messages
const (
	ErrBadStackTermination  = "invalid stack end-state: %d"
	ErrBadBranchTermination = "branches should end in the same state"
)

func verifyStackSize(code isa.Instructions) {
	s := &stackSizes{}
	s.calculateNode(visitor.Branch(code))
	if s.endSize != 0 {
		panic(fmt.Errorf(ErrBadStackTermination, s.endSize))
	}
}

// CalculateStackSize returns the maximum and final depths for the stack
// based on the instructions provided. If the final depth is non-zero,
// this is usually an indication that bad instructions were encoded
func CalculateStackSize(code isa.Instructions) (int, int) {
	s := &stackSizes{}
	s.calculateNode(visitor.Branch(code))
	return s.maxSize, s.endSize
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

func (s *stackSizes) calculateInstruction(inst *isa.Instruction) {
	oc := inst.Opcode
	effect := isa.MustGetEffect(oc)
	dPop := getStackChange(inst, effect.DPop)
	dPush := getStackChange(inst, effect.DPush)
	s.endSize += (effect.Push - effect.Pop) + (dPush - dPop)
	s.maxSize = util.IntMax(s.endSize, s.maxSize)
}

func (s *stackSizes) calculateBranches(thenNode, elseNode visitor.Node) {
	thenRes := s.calculateBranch(thenNode)
	elseRes := s.calculateBranch(elseNode)
	if elseRes.endSize != thenRes.endSize {
		panic(ErrBadBranchTermination)
	}
	s.endSize += elseRes.endSize
}

func (s *stackSizes) calculateBranch(n visitor.Node) *stackSizes {
	res := &stackSizes{
		maxSize: s.maxSize,
		endSize: 0,
	}
	res.calculateNode(n)
	s.maxSize = util.IntMax(s.maxSize, res.maxSize)
	return res
}

func getStackChange(inst *isa.Instruction, countIndex int) int {
	if countIndex > 0 {
		return int(inst.Args[countIndex-1])
	}
	return 0
}
