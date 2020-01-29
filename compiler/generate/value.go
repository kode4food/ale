package generate

import (
	"fmt"

	"github.com/kode4food/ale/compiler/encoder"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/macro"
	"github.com/kode4food/ale/namespace"
)

// Error messages
const (
	errUnknownValueType = "unknown value type: %s"
)

var consSym = namespace.RootSymbol("cons")

// Value encodes an expression
func Value(e encoder.Type, v data.Value) {
	ns := e.Globals()
	expanded := macro.Expand(ns, v)
	switch typed := expanded.(type) {
	case data.Sequence:
		Sequence(e, typed)
	case data.Pair:
		Pair(e, typed)
	case data.Symbol:
		ReferenceSymbol(e, typed)
	case data.Keyword, data.Number, data.Bool, data.Function:
		Literal(e, typed)
	default:
		panic(fmt.Errorf(errUnknownValueType, v))
	}
}

// Pair encodes a pair
func Pair(e encoder.Type, c data.Pair) {
	f := resolveBuiltIn(e, consSym)
	args := data.Values{c.Car(), c.Cdr()}
	callApplicative(e, f.Call(), args)
}

func resolveBuiltIn(e encoder.Type, sym data.Symbol) data.Caller {
	manager := e.Globals().Manager()
	root := manager.GetRoot()
	res := namespace.MustResolveValue(root, sym)
	return res.(data.Caller)
}
