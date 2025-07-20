package encoder

import (
	"github.com/kode4food/ale"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/runtime/isa"
)

// AddConstant adds a value to the constant list (if necessary)
func (e *encoder) AddConstant(val ale.Value) isa.Operand {
	c, idx := addConstant(e.constants, val)
	e.constants = c
	return idx
}

func addConstant(c data.Vector, val ale.Value) (data.Vector, isa.Operand) {
	if idx, ok := c.IndexOf(val); ok {
		return c, isa.Operand(idx)
	}
	c = append(c, val)
	return c, isa.Operand(len(c) - 1)
}
