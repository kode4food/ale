package generate

import (
	"fmt"

	"gitlab.com/kode4food/ale/api"
	"gitlab.com/kode4food/ale/internal/compiler/encoder"
	"gitlab.com/kode4food/ale/internal/macro"
)

// Error messages
const (
	UnknownValueType = "unknown value type: %s"
)

// Value encodes an expression
func Value(e encoder.Type, v api.Value) {
	ns := e.Globals()
	expanded := macro.Expand(ns, v)
	switch typed := expanded.(type) {
	case api.Sequence:
		Sequence(e, typed)
	case api.Symbol:
		Symbol(e, typed)
	case api.Keyword, api.Number, api.Bool, api.Caller, api.NilType:
		Literal(e, typed)
	default:
		panic(fmt.Errorf(UnknownValueType, v))
	}
}
