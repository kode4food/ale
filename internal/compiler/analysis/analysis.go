package analysis

import (
	"fmt"

	"gitlab.com/kode4food/ale/internal/runtime/isa"
)

type stackSizes struct {
	maxSize int
	endSize int
}

// Verify checks an ISA code stream for validity. Specifically it will
// check that jumps do not target offsets outside of the instructions
// provided and that the stack is left in a consistent state upon exit
func Verify(code isa.Instructions) {
	verifyJumps(code)
	verifyStackSize(code)
}

func verifyJumps(code isa.Instructions) {
	for _, l := range code {
		oc := l.Opcode
		if oc == isa.CondJump || oc == isa.Jump {
			mustFindLabel(code, isa.Index(l.Args[0]))
		}
	}
}

func verifyStackSize(code isa.Instructions) {
	s := &stackSizes{}
	s.calculateStackSize(code)
	if s.endSize != 0 {
		panic(fmt.Sprintf("invalid stack end-state: %d", s.endSize))
	}
}

// CalculateStackSize returns the maximum and final depths for the stack
// based on the instructions provided. If the final depth is non-zero,
// this is usually an indication that bad instructions were encoded
func CalculateStackSize(code isa.Instructions) (int, int) {
	s := &stackSizes{}
	s.calculateStackSize(code)
	return s.maxSize, s.endSize
}

func (s *stackSizes) calculateStackSize(code isa.Instructions) {
	if len(code) > 0 {
		b := Branch(code)
		s.calculatePrologue(b.Prologue)
		s.calculateBranches(b.ThenBranch, b.ElseBranch)
		s.calculateEpilogue(b.Epilogue)
	}
}

func (s *stackSizes) calculatePrologue(code isa.Instructions) {
	for _, inst := range code {
		oc := inst.Opcode
		effect := isa.MustGetEffect(oc)
		dPop := getStackChange(inst, effect.DPop)
		dPush := getStackChange(inst, effect.DPush)
		s.endSize += (effect.Push - effect.Pop) + (dPush - dPop)
		s.maxSize = maxInt(s.endSize, s.maxSize)
	}
}

func (s *stackSizes) calculateEpilogue(code isa.Instructions) {
	res := s.calculateBranch(code)
	s.endSize += res.endSize
}

func (s *stackSizes) calculateBranches(thenCode, elseCode isa.Instructions) {
	thenRes := s.calculateBranch(thenCode)
	elseRes := s.calculateBranch(elseCode)
	if elseRes.endSize != thenRes.endSize {
		panic("branches should end in the same state")
	}
	s.endSize += elseRes.endSize
}

func (s *stackSizes) calculateBranch(code isa.Instructions) *stackSizes {
	res := &stackSizes{
		maxSize: s.maxSize,
		endSize: 0,
	}
	res.calculateStackSize(code)
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
