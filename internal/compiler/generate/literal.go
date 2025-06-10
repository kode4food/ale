package generate

import (
	"github.com/kode4food/ale/internal/compiler/encoder"
	"github.com/kode4food/ale/internal/runtime/isa"
	"github.com/kode4food/ale/pkg/data"
)

// Literal encodes a literal (constant) value
func Literal(e encoder.Encoder, v data.Value) error {
	if v == data.Null {
		return Null(e)
	}
	switch v := v.(type) {
	case data.Integer:
		return Integer(e, v)
	case data.Bool:
		return Bool(e, v)
	default:
		return constant(e, v)
	}
}

// Null encodes a Null
func Null(e encoder.Encoder) error {
	e.Emit(isa.Null)
	return nil
}

// Integer encodes an Integer
func Integer(e encoder.Encoder, n data.Integer) error {
	switch {
	case n == 0:
		e.Emit(isa.Zero)
		return nil
	case n >= 0 && isa.Operand(n) <= isa.OperandMask:
		e.Emit(isa.PosInt, isa.Operand(n))
		return nil
	case n < 0 && isa.Operand(-n) <= isa.OperandMask:
		e.Emit(isa.NegInt, isa.Operand(-n))
		return nil
	}
	return constant(e, n)
}

// Bool encodes a Bool
func Bool(e encoder.Encoder, n data.Bool) error {
	if n {
		e.Emit(isa.True)
		return nil
	}
	e.Emit(isa.False)
	return nil
}

func constant(e encoder.Encoder, v data.Value) error {
	index := e.AddConstant(v)
	e.Emit(isa.Const, index)
	return nil
}
