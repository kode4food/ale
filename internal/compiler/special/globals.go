package special

import (
	"gitlab.com/kode4food/ale/api"
	"gitlab.com/kode4food/ale/internal/compiler/arity"
	"gitlab.com/kode4food/ale/internal/compiler/generate"
	"gitlab.com/kode4food/ale/internal/runtime/isa"
)

// Declare encodes global forward-declarations
func Declare(args ...api.Value) api.Value {
	e, args := splitEncoder(args)
	arity.AssertMinimum(1, len(args))
	for _, v := range args {
		name := v.(api.LocalSymbol).Name()
		generate.Literal(e, name)
		e.Append(isa.Declare)
	}
	if len(args) == 1 {
		generate.Literal(e, args[0])
	} else {
		generate.Literal(e, api.Vector(args))
	}
	return api.Nil
}

// Bind encodes a global definition
func Bind(args ...api.Value) api.Value {
	e, args := splitEncoder(args)
	arity.AssertFixed(2, len(args))
	name := args[0].(api.LocalSymbol).Name()
	generate.Value(e, args[1])
	generate.Literal(e, name)
	e.Append(isa.Bind)
	generate.Literal(e, args[0])
	return api.Nil
}
