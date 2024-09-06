package optimize

import (
	"slices"

	"github.com/kode4food/ale/internal/compiler/encoder"
	"github.com/kode4food/ale/internal/compiler/ir/visitor"
	"github.com/kode4food/ale/internal/runtime/isa"
)

type optimizer func(*encoder.Encoded)

var optimizers = [...]struct {
	process optimizer
	repeat  bool
}{
	{process: splitReturns},
	{process: makeTailCalls},
	{process: inlineCalls},
	{process: ineffectiveStores, repeat: true},
	{process: ineffectivePushes, repeat: true},
	{process: literalReturns},
}

// Encoded takes an Encoded representation and returns an optimized one
func Encoded(e *encoder.Encoded) *encoder.Encoded {
	res := e.Copy()
	for _, o := range optimizers {
	repeat:
		prev := res.Code
		o.process(res)
		if o.repeat && !slices.Equal(prev, res.Code) {
			goto repeat
		}
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
