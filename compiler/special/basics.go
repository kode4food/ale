package special

import (
	"fmt"

	"gitlab.com/kode4food/ale/api"
	"gitlab.com/kode4food/ale/compiler/arity"
	"gitlab.com/kode4food/ale/compiler/build"
	"gitlab.com/kode4food/ale/compiler/encoder"
	"gitlab.com/kode4food/ale/compiler/generate"
	"gitlab.com/kode4food/ale/eval"
	"gitlab.com/kode4food/ale/namespace"
	"gitlab.com/kode4food/ale/runtime/isa"
)

// Error messages
const (
	UnpairedBindings = "bindings must be paired"
)

// Eval encodes an evaluation
func Eval(e encoder.Type, args ...api.Value) {
	arity.AssertFixed(1, len(args))
	generate.Value(e, args[0])
	generate.Literal(e, evalFor(e.Globals()))
	e.Append(isa.Call1)
}

func evalFor(ns namespace.Type) api.Call {
	return api.Call(func(args ...api.Value) api.Value {
		return eval.Value(ns, args[0])
	})
}

// Do encodes a set of expressions, returning only the final evaluation
func Do(e encoder.Type, args ...api.Value) {
	generate.Block(e, api.Vector(args))
}

// If encodes an (if cond then else) form
func If(e encoder.Type, args ...api.Value) {
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
}

// Let encodes a (let [bindings] & body) form
func Let(e encoder.Type, args ...api.Value) {
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
}
