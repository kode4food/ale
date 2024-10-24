package generate

import (
	"github.com/kode4food/ale/internal/compiler/encoder"
	"github.com/kode4food/ale/internal/runtime/isa"
	"github.com/kode4food/ale/pkg/data"
)

// Literal encodes a literal (constant) value
func Literal(e encoder.Encoder, v data.Value) error {
	if v == data.Null {
		Null(e)
		return nil
	}
	switch v := v.(type) {
	case data.Integer, data.Float:
		Number(e, v)
	case data.Bool:
		Bool(e, v)
	default:
		index := e.AddConstant(v)
		e.Emit(isa.Const, index)
	}
	return nil
}

// Null encodes a Null
func Null(e encoder.Encoder) {
	e.Emit(isa.Null)
}

// Number encodes an Integer or Float
func Number(e encoder.Encoder, n data.Value) {
	if n, ok := n.(data.Integer); ok {
		switch {
		case n == 0:
			e.Emit(isa.Zero)
			return
		case n >= 0 && isa.Operand(n) <= isa.OperandMask:
			e.Emit(isa.PosInt, isa.Operand(n))
			return
		case n < 0 && isa.Operand(-n) <= isa.OperandMask:
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
