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
	case data.Integer:
		Integer(e, v)
	case data.Bool:
		Bool(e, v)
	default:
		Constant(e, v)
	}
	return nil
}

// Null encodes a Null
func Null(e encoder.Encoder) {
	e.Emit(isa.Null)
}

// Integer encodes an Integer
func Integer(e encoder.Encoder, n data.Integer) {
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
	Constant(e, n)
}

// Bool encodes a Bool
func Bool(e encoder.Encoder, n data.Bool) {
	if n {
		e.Emit(isa.True)
		return
	}
	e.Emit(isa.False)
}

func Constant(e encoder.Encoder, v data.Value) {
	index := e.AddConstant(v)
	e.Emit(isa.Const, index)
}
