package assert

import "github.com/kode4food/ale/runtime/isa"

// Error messages
const (
	errNonMatchingInstruction = "instruction mismatch at index %d"
)

// Instructions tests that two sets of Instructions are identical
func (w *Wrapper) Instructions(expected, actual isa.Instructions) {
	w.Helper()
	w.Equal(len(expected), len(actual))
	for i, l := range expected {
		w.Assertions.Equal(l, actual[i], errNonMatchingInstruction, i)
	}
}
