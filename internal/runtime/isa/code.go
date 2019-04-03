package isa

import "fmt"

type (
	// Code represents the atomic unit of the ISA's Code stream
	Code uint

	// Index represents a lookup offset for value arrays
	Index Code

	// Count represents a count of values
	Count Code

	// Offset represents a relative program counter offset for jumps
	Offset Code

	// Coder allows a value to return an ISA Code
	Coder interface {
		Code() Code
	}
)

// Code makes Code a Coder
func (c Code) Code() Code {
	return c
}

// Code makes Index a Coder
func (i Index) Code() Code {
	return Code(i)
}

func (i Index) String() string {
	return fmt.Sprintf("index(%d)", i)
}

// Code makes Count a Coder
func (c Count) Code() Code {
	return Code(c)
}

func (c Count) String() string {
	return fmt.Sprintf("count(%d)", c)
}

// Code makes Offset a Coder
func (o Offset) Code() Code {
	return Code(o)
}

func (o Offset) String() string {
	return fmt.Sprintf("offset(%d)", o)
}
