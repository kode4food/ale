package generate

import (
	"github.com/kode4food/ale/compiler/encoder"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/runtime/isa"
)

// Literal encodes a literal (constant) value
func Literal(e encoder.Encoder, v data.Value) {
	switch typed := v.(type) {
	case data.Null:
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
func Nil(e encoder.Encoder) {
	e.Emit(isa.Nil)
}

// Number encodes an Integer or Float
func Number(e encoder.Encoder, n data.Value) {
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
func Bool(e encoder.Encoder, n data.Bool) {
	if n {
		e.Emit(isa.True)
	} else {
		e.Emit(isa.False)
	}
}
