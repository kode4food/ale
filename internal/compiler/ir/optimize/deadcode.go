package optimize

import (
	"github.com/kode4food/ale/internal/compiler/ir/visitor"
	"github.com/kode4food/ale/internal/runtime/isa"
)

var ineffectivePushes = globalReplace(
	visitor.Pattern{
		selectEffects(func(e *isa.Effect) bool {
			return e.Push == 1 && e.Pop == 0 && !e.DPop && !e.Exit
		}),
		{isa.Pop},
	},
	removeInstructions,
)
