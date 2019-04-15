package special

import (
	"gitlab.com/kode4food/ale/api"
	"gitlab.com/kode4food/ale/internal/compiler/arity"
	"gitlab.com/kode4food/ale/internal/compiler/generate"
	"gitlab.com/kode4food/ale/internal/macro"
	"gitlab.com/kode4food/ale/internal/namespace"
	"gitlab.com/kode4food/ale/internal/runtime/isa"
)

// Quote converts its argument into a literal value
func Quote(args ...api.Value) api.Value {
	e, args := splitEncoder(args)
	arity.AssertFixed(1, len(args))
	generate.Literal(e, args[0])
	return api.Nil
}

// MacroExpand performs macro expansion of a form until it can no longer
func MacroExpand(args ...api.Value) api.Value {
	e, args := splitEncoder(args)
	arity.AssertFixed(1, len(args))
	generate.Value(e, args[0])
	generate.Literal(e, expandFor(e.Globals()))
	e.Append(isa.Call1)
	return api.Nil
}

func expandFor(ns namespace.Type) api.Call {
	return api.Call(func(args ...api.Value) api.Value {
		return macro.Expand(ns, args[0])
	})
}

// MacroExpand1 performs a single-step macro expansion of a form
func MacroExpand1(args ...api.Value) api.Value {
	e, args := splitEncoder(args)
	arity.AssertFixed(1, len(args))
	generate.Value(e, args[0])
	generate.Literal(e, expand1For(e.Globals()))
	e.Append(isa.Call1)
	return api.Nil
}

func expand1For(ns namespace.Type) api.Call {
	return api.Call(func(args ...api.Value) api.Value {
		return macro.Expand1(ns, args[0])
	})
}
