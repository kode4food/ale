package analysis

import (
	"fmt"

	"gitlab.com/kode4food/ale/internal/runtime/isa"
)

type stackSizes struct {
	maxSize int
	endSize int
}

func verifyStackSize(code isa.Instructions) {
	s := &stackSizes{}
	s.calculateNode(Branch(code))
	if s.endSize != 0 {
		panic(fmt.Sprintf("invalid stack end-state: %d", s.endSize))
	}
}

// CalculateStackSize returns the maximum and final depths for the stack
// based on the instructions provided. If the final depth is non-zero,
// this is usually an indication that bad instructions were encoded
func CalculateStackSize(code isa.Instructions) (int, int) {
	s := &stackSizes{}
	s.calculateNode(Branch(code))
	return s.maxSize, s.endSize
}

func (s *stackSizes) calculateNode(n Node) {
	switch typed := n.(type) {
	case Branches:
		s.calculateInstructions(typed.Prologue())
		s.calculateBranches(typed.ThenBranch(), typed.ElseBranch())
		s.calculateNode(typed.Epilogue())
	case Instructions:
		s.calculateInstructions(typed)
	}
}

func (s *stackSizes) calculateInstructions(inst Instructions) {
	for _, inst := range inst.Code() {
		oc := inst.Opcode
		effect := isa.MustGetEffect(oc)
		dPop := getStackChange(inst, effect.DPop)
		dPush := getStackChange(inst, effect.DPush)
		s.endSize += (effect.Push - effect.Pop) + (dPush - dPop)
		s.maxSize = maxInt(s.endSize, s.maxSize)
	}
}

func (s *stackSizes) calculateBranches(thenNode, elseNode Node) {
	thenRes := s.calculateBranch(thenNode)
	elseRes := s.calculateBranch(elseNode)
	if elseRes.endSize != thenRes.endSize {
		panic("branches should end in the same state")
	}
	s.endSize += elseRes.endSize
}

func (s *stackSizes) calculateBranch(n Node) *stackSizes {
	res := &stackSizes{
		maxSize: s.maxSize,
		endSize: 0,
	}
	res.calculateNode(n)
	s.maxSize = maxInt(s.maxSize, res.maxSize)
	return res
}

func getStackChange(inst *isa.Instruction, count int) int {
	if count > 0 {
		return int(inst.Args[count-1])
	}
	return 0
}

func maxInt(left, right int) int {
	if left > right {
		return left
	}
	return right
}
