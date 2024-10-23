package generate

import (
	"github.com/kode4food/ale/internal/compiler/encoder"
	"github.com/kode4food/ale/internal/debug"
	"github.com/kode4food/ale/pkg/data"
	"github.com/kode4food/ale/pkg/env"
	"github.com/kode4food/ale/pkg/macro"
)

var consSym = env.RootSymbol("cons")

// Value encodes an expression
func Value(e encoder.Encoder, v data.Value) error {
	ns := e.Globals()
	expanded, err := macro.Expand(ns, v)
	if err != nil {
		return err
	}
	switch expanded := expanded.(type) {
	case data.Sequence:
		return Sequence(e, expanded)
	case data.Pair:
		return Pair(e, expanded)
	case data.Symbol:
		return ReferenceSymbol(e, expanded)
	case data.Keyword, data.Number, data.Bool, data.Procedure:
		return Literal(e, expanded)
	default:
		panic(debug.ProgrammerError("unknown value type: %s", v))
	}
}

// Pair encodes a pair
func Pair(e encoder.Encoder, c data.Pair) error {
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
