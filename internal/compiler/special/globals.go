package special

import (
	"gitlab.com/kode4food/ale/api"
	"gitlab.com/kode4food/ale/internal/compiler/arity"
	"gitlab.com/kode4food/ale/internal/compiler/generate"
	"gitlab.com/kode4food/ale/internal/runtime/isa"
)

// Declare encodes a global forward-declaration
func Declare(args ...api.Value) api.Value {
	e, args := splitEncoder(args)
	arity.AssertFixed(1, len(args))
	name := args[0].(api.LocalSymbol).Name()
	generate.Literal(e, name)
	e.Append(isa.Declare)
	generate.Literal(e, name)
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
	generate.Literal(e, name)
	return api.Nil
}
