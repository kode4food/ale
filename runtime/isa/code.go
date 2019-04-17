package isa

import "fmt"

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

	// Instructions represents a set of Instructions
	Instructions []*Instruction
)

// Error messages
const (
	BadInstructionArgs = "instruction argument mismatch: %s"
)

// Word makes Index a Coder
func (i Index) Word() Word {
	return Word(i)
}

func (i Index) String() string {
	return fmt.Sprintf("index(%d)", i)
}

// Word makes Count a Coder
func (c Count) Word() Word {
	return Word(c)
}

func (c Count) String() string {
	return fmt.Sprintf("count(%d)", c)
}

// Word makes Offset a Coder
func (o Offset) Word() Word {
	return Word(o)
}

func (o Offset) String() string {
	return fmt.Sprintf("offset(%d)", o)
}

// New creates a new Instruction instance
func New(oc Opcode, args ...Word) *Instruction {
	effect := MustGetEffect(oc)
	if len(args) != effect.Size-1 {
		panic(fmt.Sprintf(BadInstructionArgs, oc))
	}
	return &Instruction{
		Opcode: oc,
		Args:   args,
	}
}
