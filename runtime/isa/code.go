package isa

import (
	"fmt"
	"strings"

	"github.com/kode4food/ale/data"
	str "github.com/kode4food/ale/internal/strings"
	"github.com/kode4food/comb/slices"
)

type (
	// Word represents the atomic unit of the ISA's Instruction stream. We've
	// chosen a size that best aligns to the host architecture's standard word
	// size for alignment reasons
	Word uintptr

	// Instruction represents a single instruction and its operand
	Instruction Word

	// Instructions represent a set of Instructions
	Instructions []Instruction
)

// Error messages
const (
	ErrBadInstruction  = "instruction operand mismatch: %s"
	ErrExpectedOperand = "expected unsigned operand: %d"
)

// New creates a new Instruction instance
func New(oc Opcode, args ...Operand) Instruction {
	effect := MustGetEffect(oc)
	switch {
	case effect.Operand != Nothing && len(args) == 1:
		if !IsValidOperand(int(args[0])) {
			panic(fmt.Errorf(ErrExpectedOperand, args[0]))
		}
		return Instruction(Opcode(args[0]<<OpcodeSize) | oc)
	case effect.Operand == Nothing && len(args) == 0:
		return Instruction(oc)
	default:
		panic(fmt.Errorf(ErrBadInstruction, oc.String()))
	}
}

func (i Instructions) String() string {
	return strings.Join(
		slices.Map(i, func(in Instruction) string {
			return in.String()
		}),
		"\n",
	)
}

func (i Instruction) Split() (Opcode, Operand) {
	return Opcode(i) & OpcodeMask, Operand(i) >> OpcodeSize & OperandMask
}

// Equal compares this Instruction to another for equality
func (i Instruction) Equal(v data.Value) bool {
	if v, ok := v.(Instruction); ok {
		return i == v
	}
	return false
}

func (i Instruction) String() string {
	oc, operand := i.Split()
	s := str.CamelToSnake(oc.String())
	if Effects[oc].Operand != Nothing {
		return fmt.Sprintf("%s %d", s, operand)
	}
	return s
}

// IsValidOperand returns true if the int falls within the operand range
func IsValidOperand(i int) bool {
	return i >= 0 && Operand(i) <= OperandMask
}
