package analysis

import "github.com/kode4food/ale/runtime/isa"

// Verify checks an ISA code stream for validity. Specifically it will check
// that jumps do not target offsets outside the instructions provided and that
// the stack is left in a consistent state upon exit
func Verify(code isa.Instructions) {
	verifyJumps(code)
	verifyStackSize(code)
}
