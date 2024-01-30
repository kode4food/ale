package optimize

import (
	"github.com/kode4food/ale/compiler/encoder"
	"github.com/kode4food/ale/compiler/ir/visitor"
	"github.com/kode4food/comb/basics"
)

type optimizer func(visitor.Node) visitor.Node

var makeOptimizers = []func(*encoder.Encoded) optimizer{
	makeSplitReturns,   // roll standalone returns into preceding branches
	makeTailCalls,      // replace calls in tail position with a tail-call
	makeLiteralReturns, // convert literal returns into single instructions
}

// FromEncoded takes an Encoded representation and returns an optimized one
func FromEncoded(e *encoder.Encoded) *encoder.Encoded {
	res := *e
	res.Code = basics.FoldLeft(
		basics.Map(makeOptimizers,
			func(fn func(*encoder.Encoded) optimizer) optimizer {
				return fn(&res)
			},
		),
		visitor.Branch(res.Code),
		func(node visitor.Node, o optimizer) visitor.Node {
			return o(node)
		},
	).Code()
	return &res
}
