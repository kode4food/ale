package encoder

import (
	"reflect"

	"gitlab.com/kode4food/ale/data"
	"gitlab.com/kode4food/ale/runtime/isa"
)

// Constants returns the encoder's set of constants
func (e *encoder) Constants() data.Values {
	ec := e.constants
	res := make(data.Values, len(ec))
	copy(res, ec)
	return res
}

// AddConstant adds a value to the constant list (if necessary)
func (e *encoder) AddConstant(val data.Value) isa.Index {
	for i, c := range e.constants {
		if reflect.DeepEqual(c, val) {
			return isa.Index(i)
		}
	}
	c := append(e.constants, val)
	e.constants = c
	return isa.Index(len(c) - 1)
}
