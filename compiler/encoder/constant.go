package encoder

import (
	"slices"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/runtime/isa"
)

// Constants returns the encoder's set of constants
func (e *encoder) Constants() data.Values {
	ec := e.constants
	res := make(data.Values, len(ec))
	copy(res, ec)
	return res
}

// AddConstant adds a value to the constant list (if necessary)
func (e *encoder) AddConstant(val data.Value) isa.Operand {
	if idx, ok := e.findConstant(val); ok {
		return isa.Operand(idx)
	}
	c := append(e.constants, val)
	e.constants = c
	return isa.Operand(len(c) - 1)
}

func (e *encoder) findConstant(val data.Value) (int, bool) {
	i := slices.IndexFunc(e.constants, val.Equal)
	return i, i != -1
}
