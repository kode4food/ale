package generate

import (
	"github.com/kode4food/ale"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/env"
	"github.com/kode4food/ale/internal/compiler/encoder"
	"github.com/kode4food/ale/internal/runtime/isa"
	"github.com/kode4food/ale/macro"
)

// Value encodes an expression
func Value(e encoder.Encoder, v ale.Value) error {
	ns := e.Globals()
	ex, err := macro.Expand(ns, v)
	if err != nil {
		return err
	}
	return expanded(e, ex)
}

func expanded(e encoder.Encoder, v ale.Value) error {
	switch v := v.(type) {
	case data.Qualified:
		return Global(e, v)
	case data.Local:
		return Reference(e, v)
	case *data.List:
		return Call(e, v)
	case data.Vector:
		return Vector(e, v)
	case *data.Object:
		return Object(e, v)
	case *data.Cons:
		return Cons(e, v)
	default:
		return Literal(e, v)
	}
}

// Cons encodes a cons pair
func Cons(e encoder.Encoder, c *data.Cons) error {
	if err := Value(e, c.Cdr()); err != nil {
		return err
	}
	if err := Value(e, c.Car()); err != nil {
		return err
	}
	e.Emit(isa.Cons)
	return nil
}

func resolveBuiltIn(
	e encoder.Encoder, sym data.Symbol,
) (data.Procedure, error) {
	ge := e.Globals().Environment()
	root := ge.GetRoot()
	res, err := env.ResolveValue(root, sym)
	if err != nil {
		return nil, err
	}
	return res.(data.Procedure), nil
}
