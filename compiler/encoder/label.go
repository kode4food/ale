package encoder

import "github.com/kode4food/ale/runtime/isa"

// NewLabel allocates a Label (
func (e *encoder) NewLabel() isa.Operand {
	res := e.nextLabel
	e.nextLabel++
	return res
}
