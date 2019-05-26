package optimize

import (
	"gitlab.com/kode4food/ale/compiler/ir/visitor"
	"gitlab.com/kode4food/ale/runtime/isa"
)

type optimizer func(visitor.Node) visitor.Node

var optimizers = []optimizer{
	splitReturns,
	tailCalls,
	literalReturns,
}

// Instructions performs optimizations on the provided instructions
func Instructions(code isa.Instructions) isa.Instructions {
	root := visitor.Branch(code)
	for _, o := range optimizers {
		root = o(root)
	}
	return root.Code()
}
