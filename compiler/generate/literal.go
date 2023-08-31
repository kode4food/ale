package generate

import (
	"github.com/kode4food/ale/compiler/encoder"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/runtime/isa"
)

// Literal encodes a literal (constant) value
func Literal(e encoder.Encoder, v data.Value) {
	switch v := v.(type) {
	case data.Null:
		Nil(e)
	case data.Integer, data.Float:
		Number(e, v)
	case data.Bool:
		Bool(e, v)
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
	if n, ok := n.(data.Integer); ok {
		switch {
		case n == 0:
			e.Emit(isa.Zero)
			return
		case n >= 0 && n <= isa.OperandMask:
			e.Emit(isa.PosInt, isa.Operand(n))
			return
		case n < 0 && -n <= isa.OperandMask:
			e.Emit(isa.NegInt, isa.Operand(-n))
			return
		}
	}
	index := e.AddConstant(n)
	e.Emit(isa.Const, index)
}

// Bool encodes a Bool
func Bool(e encoder.Encoder, n data.Bool) {
	if n {
		e.Emit(isa.True)
	} else {
		e.Emit(isa.False)
	}
}
