package generate

import (
	"gitlab.com/kode4food/ale/api"
	"gitlab.com/kode4food/ale/internal/compiler/encoder"
	"gitlab.com/kode4food/ale/internal/runtime/isa"
)

// Literal encodes a literal (constant) value
func Literal(e encoder.Type, v api.Value) {
	switch typed := v.(type) {
	case api.NilType:
		Nil(e)
	case api.Integer, api.Float:
		Number(e, typed)
	case api.Bool:
		Bool(e, typed)
	default:
		index := e.AddConstant(v)
		e.Append(isa.Const, index)
	}
}

// Nil encodes a Nil
func Nil(e encoder.Type) {
	e.Append(isa.Nil)
}

// Number encodes an Integer or Float
func Number(e encoder.Type, n api.Value) {
	switch n {
	case api.Integer(0):
		e.Append(isa.Zero)
	case api.Integer(1):
		e.Append(isa.One)
	case api.Integer(2):
		e.Append(isa.Two)
	case api.Integer(-1):
		e.Append(isa.NegOne)
	default:
		index := e.AddConstant(n)
		e.Append(isa.Const, index)
	}
}

// Bool encodes a Bool
func Bool(e encoder.Type, n api.Bool) {
	if n {
		e.Append(isa.True)
	} else {
		e.Append(isa.False)
	}
}
