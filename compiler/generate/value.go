package generate

import (
	"fmt"

	"github.com/kode4food/ale/compiler/encoder"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/env"
	"github.com/kode4food/ale/macro"
)

// Error messages
const (
	errUnknownValueType = "unknown value type: %s"
)

var consSym = env.RootSymbol("cons")

// Value encodes an expression
func Value(e encoder.Encoder, v data.Value) {
	ns := e.Globals()
	switch expanded := macro.Expand(ns, v).(type) {
	case data.Sequence:
		Sequence(e, expanded)
	case data.Pair:
		Pair(e, expanded)
	case data.Symbol:
		ReferenceSymbol(e, expanded)
	case data.Keyword, data.Number, data.Bool, data.Function:
		Literal(e, expanded)
	default:
		// Programmer error
		panic(fmt.Errorf(errUnknownValueType, v))
	}
}

// Pair encodes a pair
func Pair(e encoder.Encoder, c data.Pair) {
	f := resolveBuiltIn(e, consSym)
	args := data.Values{c.Car(), c.Cdr()}
	callApplicative(e, f, args)
}

func resolveBuiltIn(e encoder.Encoder, sym data.Symbol) data.Function {
	ge := e.Globals().Environment()
	root := ge.GetRoot()
	res := env.MustResolveValue(root, sym)
	return res.(data.Function)
}
