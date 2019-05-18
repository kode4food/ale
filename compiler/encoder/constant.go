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
	if idx, ok := e.findConstant(val); ok {
		return isa.Index(idx)
	}
	c := append(e.constants, val)
	e.constants = c
	return isa.Index(len(c) - 1)
}

func (e *encoder) findConstant(val data.Value) (int, bool) {
	if _, ok := val.(data.Call); ok {
		return -1, false
	}
	for i, c := range e.constants {
		if _, ok := c.(data.Call); ok {
			continue
		}
		if c == val || reflect.DeepEqual(c, val) {
			return i, true
		}
	}
	return -1, false
}
