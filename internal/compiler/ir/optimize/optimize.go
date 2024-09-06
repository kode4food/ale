package optimize

import (
	"slices"

	"github.com/kode4food/ale/internal/compiler/encoder"
	"github.com/kode4food/ale/internal/compiler/ir/visitor"
	"github.com/kode4food/ale/internal/runtime/isa"
)

type optimizer func(*encoder.Encoded)

var optimizers = [...]optimizer{
	splitReturns,
	makeTailCalls,
	inlineCalls,
	repeatWhenModified(ineffectiveStores),
	repeatWhenModified(ineffectivePushes),
	literalReturns,
}

// Encoded takes an Encoded representation and returns an optimized one
func Encoded(e *encoder.Encoded) *encoder.Encoded {
	res := e.Copy()
	for _, o := range optimizers {
		o(res)
	}
	return res
}

func repeatWhenModified(o optimizer) optimizer {
	return func(e *encoder.Encoded) {
		for {
			prev := e.Code
			o(e)
			if slices.Equal(prev, e.Code) {
				return
			}
		}
	}
}

func globalReplace(p visitor.Pattern, m visitor.Mapper) optimizer {
	replace := visitor.Replace(p, m)
	return func(e *encoder.Encoded) {
		root := visitor.All(e.Code)
		visitor.Visit(root, replace)
		e.Code = root.Code()
	}
}

func selectEffects(filter func(*isa.Effect) bool) []isa.Opcode {
	var res []isa.Opcode
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
