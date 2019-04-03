package special

import (
	"fmt"

	"gitlab.com/kode4food/ale/api"
	"gitlab.com/kode4food/ale/eval"
	"gitlab.com/kode4food/ale/internal/compiler/arity"
	"gitlab.com/kode4food/ale/internal/compiler/build"
	"gitlab.com/kode4food/ale/internal/compiler/encoder"
	"gitlab.com/kode4food/ale/internal/compiler/generate"
	"gitlab.com/kode4food/ale/internal/macro"
	"gitlab.com/kode4food/ale/internal/namespace"
	"gitlab.com/kode4food/ale/internal/runtime/isa"
)

// Error messages
const (
	UnpairedBindings = "bindings must be paired"
)

func splitEncoder(args api.Values) (encoder.Type, api.Values) {
	e := args[0].(encoder.Type)
	return e, args[1:]
}

// Eval encodes an evaluation
func Eval(args ...api.Value) api.Value {
	e, args := splitEncoder(args)
	arity.AssertFixed(1, len(args))
	generate.Value(e, args[0])
	generate.Literal(e, evalFor(e.Globals()))
	e.Append(isa.Call1)
	return api.Nil
}

func evalFor(ns namespace.Type) api.Call {
	return api.Call(func(args ...api.Value) api.Value {
		return eval.Value(ns, args[0])
	})
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

// Do encodes a set of expressions, returning only the final evaluation
func Do(args ...api.Value) api.Value {
	e, args := splitEncoder(args)
	generate.Block(e, api.Vector(args))
	return api.Nil
}

// If encodes an (if cond then else) form
func If(args ...api.Value) api.Value {
	e, args := splitEncoder(args)
	al := arity.AssertRanged(2, 3, len(args))
	build.Cond(e,
		func() {
			generate.Value(e, args[0])
			e.Append(isa.MakeTruthy)
		},
		func() {
			generate.Value(e, args[1])
		},
		func() {
			if al == 3 {
				generate.Value(e, args[2])
			} else {
				generate.Nil(e)
			}
		},
	)
	return api.Nil
}

// Let encodes a (let [bindings] & body) form
func Let(args ...api.Value) api.Value {
	e, args := splitEncoder(args)
	arity.AssertMinimum(2, len(args))
	bindings := args[0].(api.Vector)
	lb := len(bindings)
	if lb%2 != 0 {
		panic(fmt.Errorf(UnpairedBindings))
	}

	for i := 0; i < lb; i += 2 {
		n := bindings[i].(api.LocalSymbol).Name()
		e.PushLocals()
		generate.Value(e, bindings[i+1])
		e.Append(isa.Store, e.AddLocal(n))
	}

	body := api.Vector(args[1:])
	generate.Block(e, body)

	for i := 0; i < lb; i += 2 {
		e.PopLocals()
	}
	return api.Nil
}
