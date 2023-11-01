package optimize

import (
	"github.com/kode4food/ale/compiler/ir/visitor"
	"github.com/kode4food/ale/runtime/isa"
	"github.com/kode4food/comb/basics"
)

type optimizer func(visitor.Node) visitor.Node

var optimizers = []optimizer{
	splitReturns,   // roll standalone returns into preceding branches
	tailCalls,      // replace calls in tail position with a tail-call
	literalReturns, // convert some literal returns into single instructions
}

// Instructions perform optimizations on the provided instructions
func Instructions(code isa.Instructions) isa.Instructions {
	return basics.FoldLeft(
		optimizers,
		visitor.Branch(code),
		func(node visitor.Node, o optimizer) visitor.Node {
			return o(node)
		},
	).Code()
}
