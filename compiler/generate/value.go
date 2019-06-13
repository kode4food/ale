package generate

import (
	"fmt"

	"gitlab.com/kode4food/ale/compiler/encoder"
	"gitlab.com/kode4food/ale/data"
	"gitlab.com/kode4food/ale/macro"
)

// Error messages
const (
	UnknownValueType = "unknown value type: %s"
)

// Value encodes an expression
func Value(e encoder.Type, v data.Value) {
	ns := e.Globals()
	expanded := macro.Expand(ns, v)
	switch typed := expanded.(type) {
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
