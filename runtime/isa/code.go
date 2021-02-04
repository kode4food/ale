package isa

import (
	"fmt"
	"strings"

	"github.com/kode4food/ale/data"
)

type (
	// Word represents the atomic unit of the ISA's code stream
	Word uint

	// Index represents a lookup offset for value arrays
	Index Word

	// Count represents a count of values
	Count Word

	// Offset represents a relative program counter offset for jumps
	Offset Word

	// Coder allows a value to return an ISA Word
	Coder interface {
		Word() Word
	}

	// Instruction represents a single instruction and its arguments
	Instruction struct {
		Opcode
		Args []Word
	}

	// Instructions represent a set of Instructions
	Instructions []*Instruction
)

// Error messages
const (
	ErrBadInstructionArgs = "instruction argument mismatch: %s"
)

// Word makes Index a Word
func (i Index) Word() Word {
	return Word(i)
}

// Word makes Count a Word
func (c Count) Word() Word {
	return Word(c)
}

// Word makes Offset a Word
func (o Offset) Word() Word {
	return Word(o)
}

// New creates a new Instruction instance
func New(oc Opcode, args ...Word) *Instruction {
	effect := MustGetEffect(oc)
	if len(args) != effect.Size-1 {
		panic(fmt.Errorf(ErrBadInstructionArgs, oc.String()))
	}
	return &Instruction{
		Opcode: oc,
		Args:   args,
	}
}

// Equal compares this Instruction to another for equality
func (i *Instruction) Equal(v data.Value) bool {
	if v, ok := v.(*Instruction); ok {
		if i.Opcode != v.Opcode || len(i.Args) != len(v.Args) {
			return false
		}
		for i, l := range i.Args {
			if l != v.Args[i] {
				return false
			}
		}
		return true
	}
	return false
}

func (i *Instruction) String() string {
	args := i.Args
	strs := make([]string, len(args))
	for i, a := range args {
		strs[i] = fmt.Sprintf("%d", a)
	}
	joined := strings.Join(strs, ", ")
	return fmt.Sprintf("%s(%s)", i.Opcode.String(), joined)
}
