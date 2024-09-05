package optimize

import (
	"slices"

	"github.com/kode4food/ale/internal/compiler/encoder"
	"github.com/kode4food/ale/internal/compiler/ir/visitor"
	"github.com/kode4food/ale/internal/runtime/isa"
)

type optimizer func(*encoder.Encoded)

var optimizers = [...]optimizer{
	splitReturns,      // roll standalone returns into preceding branches
	makeTailCalls,     // replace calls in tail position with a tail-call
	inlineCalls,       // inline calls to procedures that qualify
	ineffectiveStores, // isolated store followed by a load of same operand
	ineffectivePushes, // values pushed to the stack for no reason
	literalReturns,    // convert literal returns into single instructions
}

// Encoded takes an Encoded representation and returns an optimized one
func Encoded(e *encoder.Encoded) *encoder.Encoded {
	res := e.Copy()
	for _, o := range optimizers {
		o(res)
	}
	return res
}

func globalReplace(p visitor.Pattern, m visitor.Mapper) optimizer {
	replace := visitor.Replace(p, m)
	return func(e *encoder.Encoded) {
		root := visitor.All(e.Code)
		visitor.Visit(root, replace)
		e.Code = root.Code()
	}
}

func globalRepeatedReplace(p visitor.Pattern, m visitor.Mapper) optimizer {
	replace := visitor.Replace(p, m)
	return func(e *encoder.Encoded) {
		for {
			root := visitor.All(e.Code)
			visitor.Visit(root, replace)
			res := root.Code()
			if slices.Equal(e.Code, res) {
				return
			}
			e.Code = res
		}
	}
}

func selectEffects(filter func(*isa.Effect) bool) []isa.Opcode {
	res := make([]isa.Opcode, 0)
	for oc, effect := range isa.Effects {
		if filter(effect) {
			res = append(res, oc)
		}
	}
	return res
}

func removeInstructions(_ isa.Instructions) isa.Instructions {
	return isa.Instructions{}
}
