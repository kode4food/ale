package generate

import (
	"github.com/kode4food/ale/internal/compiler/encoder"
	"github.com/kode4food/ale/pkg/data"
	"github.com/kode4food/ale/pkg/env"
	"github.com/kode4food/ale/pkg/macro"
)

var consSym = env.RootSymbol("cons")

// Value encodes an expression
func Value(e encoder.Encoder, v data.Value) error {
	ns := e.Globals()
	ex, err := macro.Expand(ns, v)
	if err != nil {
		return err
	}
	return expanded(e, ex)
}

func expanded(e encoder.Encoder, v data.Value) error {
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
	f, err := resolveBuiltIn(e, consSym)
	if err != nil {
		return err
	}
	args := data.Vector{c.Car(), c.Cdr()}
	return callStatic(e, f, args)
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
