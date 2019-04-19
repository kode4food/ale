package generate

import (
	"gitlab.com/kode4food/ale/compiler/encoder"
	"gitlab.com/kode4food/ale/data"
	"gitlab.com/kode4food/ale/runtime/isa"
)

// Literal encodes a literal (constant) value
func Literal(e encoder.Type, v data.Value) {
	switch typed := v.(type) {
	case data.NilType:
		Nil(e)
	case data.Integer, data.Float:
		Number(e, typed)
	case data.Bool:
		Bool(e, typed)
	default:
		index := e.AddConstant(v)
		e.Emit(isa.Const, index)
	}
}

// Nil encodes a Nil
func Nil(e encoder.Type) {
	e.Emit(isa.Nil)
}

// Number encodes an Integer or Float
func Number(e encoder.Type, n data.Value) {
	switch n {
	case data.Integer(0):
		e.Emit(isa.Zero)
	case data.Integer(1):
		e.Emit(isa.One)
	case data.Integer(2):
		e.Emit(isa.Two)
	case data.Integer(-1):
		e.Emit(isa.NegOne)
	default:
		index := e.AddConstant(n)
		e.Emit(isa.Const, index)
	}
}

// Bool encodes a Bool
func Bool(e encoder.Type, n data.Bool) {
	if n {
		e.Emit(isa.True)
	} else {
		e.Emit(isa.False)
	}
}
