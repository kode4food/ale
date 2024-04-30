package analysis

import "github.com/kode4food/ale/internal/runtime/isa"

// Verify checks an ISA Instruction stream for validity. Specifically, it will
// check that jumps do not target offsets outside the instructions provided and
// that the stack is left in a consistent state upon exit
func Verify(code isa.Instructions) error {
	if err := verifyJumps(code); err != nil {
		return err
	}
	if err := verifyStackSize(code); err != nil {
		return err
	}
	return nil
}

// MustVerify is a wrapper around Verify that will panic if the Instructions
// provided are invalid
func MustVerify(code isa.Instructions) {
	if err := Verify(code); err != nil {
		panic(err)
	}
}
