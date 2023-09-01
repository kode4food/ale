package special

import (
	"github.com/kode4food/ale/compiler/encoder"
	"github.com/kode4food/ale/compiler/generate"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/env"
	"github.com/kode4food/ale/macro"
	"github.com/kode4food/ale/runtime/isa"
)

// MacroExpand performs macro expansion of a form until it can no longer
func MacroExpand(e encoder.Encoder, args ...data.Value) {
	data.AssertFixed(1, len(args))
	generate.Value(e, args[0])
	generate.Literal(e, expandFor(e.Globals()))
	e.Emit(isa.Call1)
}

func expandFor(ns env.Namespace) data.Function {
	return data.Applicative(func(args ...data.Value) data.Value {
		return macro.Expand(ns, args[0])
	}, 1)
}

// MacroExpand1 performs a single-step macro expansion of a form
func MacroExpand1(e encoder.Encoder, args ...data.Value) {
	data.AssertFixed(1, len(args))
	generate.Value(e, args[0])
	generate.Literal(e, expand1For(e.Globals()))
	e.Emit(isa.Call1)
}

func expand1For(ns env.Namespace) data.Function {
	return data.Applicative(func(args ...data.Value) data.Value {
		return macro.Expand1(ns, args[0])
	}, 1)
}
