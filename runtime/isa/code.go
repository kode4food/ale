package isa

import (
	"fmt"
	"math"
	"strings"

	"github.com/kode4food/ale/data"
)

type (
	// Word represents the atomic unit of the ISA's code stream
	Word uint32

	// Instruction represents a single instruction and its operand
	Instruction Word

	Coder interface {
		Instruction() Instruction
	}

	// Operand parameterizes an Instruction
	Operand Word

	// Instructions represent a set of Instructions
	Instructions []Instruction
)

// MaxWord is the highest value of an Instruction Word
const MaxWord = math.MaxUint32

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

func (i Instruction) Instruction() Instruction {
	return i
}

func (i Instructions) String() string {
	strs := make([]string, len(i))
	for j, l := range i {
		strs[j] = l.String()
	}
	return strings.Join(strs, "\n")
}

func (i Instruction) Split() (Opcode, Operand) {
	return Opcode(i & OpcodeMask), Operand(i >> OpcodeSize & OperandMask)
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
	if Effects[oc].Operand != Nothing {
		return fmt.Sprintf("%s(%d)", oc.String(), operand)
	}
	return fmt.Sprintf("%s()", oc.String())
}

// IsValidOperand returns true if the int falls within the operand range
func IsValidOperand(i int) bool {
	return i >= 0 && i <= OperandMask
}
