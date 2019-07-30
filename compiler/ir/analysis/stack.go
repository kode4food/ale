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

func verifyStackSize(code isa.Instructions) {
	s := &stackSizes{}
	s.calculateNode(visitor.Branch(code))
	if s.endSize != 0 {
		panic(fmt.Sprintf("invalid stack end-state: %d", s.endSize))
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
	switch typed := n.(type) {
	case visitor.Branches:
		s.calculateInstructions(typed.Prologue())
		s.calculateBranches(typed.ThenBranch(), typed.ElseBranch())
		s.calculateNode(typed.Epilogue())
	case visitor.Instructions:
		s.calculateInstructions(typed)
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
		panic("branches should end in the same state")
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

func getStackChange(inst *isa.Instruction, count int) int {
	if count > 0 {
		return int(inst.Args[count-1])
	}
	return 0
}
