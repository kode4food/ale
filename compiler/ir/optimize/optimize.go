package optimize

import (
	"github.com/kode4food/ale/compiler/encoder"
	"github.com/kode4food/ale/runtime/isa"
	"github.com/kode4food/comb/basics"
)

type optimizer func(instructions isa.Instructions) isa.Instructions

var makeOptimizers = []func(*encoder.Encoded) optimizer{
	makeSplitReturns,   // roll standalone returns into preceding branches
	makeLiteralReturns, // convert literal returns into single instructions
	makeTailCalls,      // replace calls in tail position with a tail-call
	makeInlineCalls,    // inline calls to procedures that qualify
}

// Encoded takes an Encoded representation and returns an optimized one
func Encoded(e *encoder.Encoded) *encoder.Encoded {
	res := *e
	optimizers := basics.Map(makeOptimizers,
		func(fn func(*encoder.Encoded) optimizer) optimizer {
			return fn(&res)
		},
	)
	for _, o := range optimizers {
		res.Code = o(res.Code)
	}
	return &res
}
