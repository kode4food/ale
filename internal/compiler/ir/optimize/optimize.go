package optimize

import (
	"github.com/kode4food/ale/internal/basics"
	"github.com/kode4food/ale/internal/compiler/encoder"
	"github.com/kode4food/ale/internal/compiler/ir/visitor"
	"github.com/kode4food/ale/internal/runtime/isa"
)

type optimizer func(*encoder.Encoded) *encoder.Encoded

// Encoded takes an Encoded representation and returns an optimized one
var Encoded = compose(
	splitReturns,
	callsInTailPosition,
	inlineCalls,
	repeatWhenModified(redundantLocals, ineffectivePushes),
	literalReturns,
	dedicatedCalls,
)

func compose(first optimizer, rest ...optimizer) optimizer {
	if len(rest) == 0 {
		return first
	}
	optimizers := append([]optimizer{first}, rest...)
	return func(e *encoder.Encoded) *encoder.Encoded {
		res := e
		for _, o := range optimizers {
			res = o(res)
		}
		return res
	}
}

func repeatWhenModified(first optimizer, rest ...optimizer) optimizer {
	o := compose(first, rest...)
	return func(e *encoder.Encoded) *encoder.Encoded {
		res := e
		for {
			prev := res.Code
			res = o(res)
			if basics.Equal(prev, res.Code) {
				return res
			}
		}
	}
}

func globalReplace(p visitor.Pattern, m visitor.Mapper) optimizer {
	return func(e *encoder.Encoded) *encoder.Encoded {
		r := visitor.Replace(p, m)
		return performReplace(e, r)
	}
}

func performReplace(e *encoder.Encoded, r *visitor.Replacer) *encoder.Encoded {
	root := visitor.All(e.Code)
	visitor.Visit(root, r)
	return e.WithCode(root.Code())
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
