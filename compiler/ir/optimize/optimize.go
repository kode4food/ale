package optimize

import (
	"github.com/kode4food/ale/compiler/ir/visitor"
	"github.com/kode4food/ale/runtime/isa"
)

type optimizer func(visitor.Node) visitor.Node

var optimizers = []optimizer{
	splitReturns,   // roll standalone returns into preceding branches
	tailCalls,      // replace calls in tail position with a tail-call
	literalReturns, // convert some literal returns into single instructions
	unTruthy,       // instructions that push bool don't need MakeTruthy
}

// Instructions performs optimizations on the provided instructions
func Instructions(code isa.Instructions) isa.Instructions {
	root := visitor.Branch(code)
	for _, o := range optimizers {
		root = o(root)
	}
	return root.Code()
}
