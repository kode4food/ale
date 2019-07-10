package generate

import (
	"fmt"

	"gitlab.com/kode4food/ale/compiler/encoder"
	"gitlab.com/kode4food/ale/data"
	"gitlab.com/kode4food/ale/macro"
	"gitlab.com/kode4food/ale/namespace"
)

// Error messages
const (
	UnknownValueType = "unknown value type: %s"
)

var consSym = namespace.RootSymbol("cons")

// Value encodes an expression
func Value(e encoder.Type, v data.Value) {
	ns := e.Globals()
	expanded := macro.Expand(ns, v)
	switch typed := expanded.(type) {
	case *data.Cons:
		Pair(e, typed)
	case data.Sequence:
		Sequence(e, typed)
	case data.Symbol:
		ReferenceSymbol(e, typed)
	case data.Keyword, data.Number, data.Bool, data.Function:
		Literal(e, typed)
	default:
		panic(fmt.Errorf(UnknownValueType, v))
	}
}

// Pair encodes a Cons pair
func Pair(e encoder.Type, c *data.Cons) {
	f := resolveBuiltIn(e, consSym)
	args := data.Values{c.Car(), c.Cdr()}
	callApplicative(e, f.Caller(), args)
}

func resolveBuiltIn(e encoder.Type, sym data.Symbol) data.Caller {
	manager := e.Globals().Manager()
	root := manager.GetRoot()
	res := namespace.MustResolveValue(root, sym)
	return res.(data.Caller)
}
