package optimize

import (
	"github.com/kode4food/ale/internal/compiler/ir/visitor"
	"github.com/kode4food/ale/internal/runtime/isa"
)

var (
	// dedicatedCalls convert generic Call instructions into dedicated ones
	dedicatedCalls = globalReplace(
		visitor.Pattern{{isa.Call}},
		dedicatedCallMapper,
	)

	mappedDedicatedCalls = map[isa.Operand]isa.Instructions{
		0: {isa.Call0.New()},
		1: {isa.Call1.New()},
		2: {isa.Call2.New()},
		3: {isa.Call3.New()},
	}
)

func dedicatedCallMapper(i isa.Instructions) isa.Instructions {
	if res, ok := mappedDedicatedCalls[i[0].Operand()]; ok {
		return res
	}
	return i
}
