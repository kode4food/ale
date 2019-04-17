package special

import (
	"gitlab.com/kode4food/ale/api"
	"gitlab.com/kode4food/ale/compiler/arity"
	"gitlab.com/kode4food/ale/compiler/encoder"
	"gitlab.com/kode4food/ale/compiler/generate"
	"gitlab.com/kode4food/ale/macro"
	"gitlab.com/kode4food/ale/namespace"
	"gitlab.com/kode4food/ale/runtime/isa"
)

// Quote converts its argument into a literal value
func Quote(e encoder.Type, args ...api.Value) {
	arity.AssertFixed(1, len(args))
	generate.Literal(e, args[0])
}

// MacroExpand performs macro expansion of a form until it can no longer
func MacroExpand(e encoder.Type, args ...api.Value) {
	arity.AssertFixed(1, len(args))
	generate.Value(e, args[0])
	generate.Literal(e, expandFor(e.Globals()))
	e.Append(isa.Call1)
}

func expandFor(ns namespace.Type) api.Call {
	return api.Call(func(args ...api.Value) api.Value {
		return macro.Expand(ns, args[0])
	})
}

// MacroExpand1 performs a single-step macro expansion of a form
func MacroExpand1(e encoder.Type, args ...api.Value) {
	arity.AssertFixed(1, len(args))
	generate.Value(e, args[0])
	generate.Literal(e, expand1For(e.Globals()))
	e.Append(isa.Call1)
}

func expand1For(ns namespace.Type) api.Call {
	return api.Call(func(args ...api.Value) api.Value {
		return macro.Expand1(ns, args[0])
	})
}
