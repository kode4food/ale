package optimize

import (
	"github.com/kode4food/ale/internal/basics"
	"github.com/kode4food/ale/internal/compiler/ir/visitor"
	"github.com/kode4food/ale/internal/runtime/isa"
)

var (
	literalReturnMap = map[isa.Opcode]isa.Opcode{
		isa.False: isa.RetFalse,
		isa.Null:  isa.RetNull,
		isa.True:  isa.RetTrue,
	}

	// literalReturns convert some literal instructions followed by a return
	// instruction into single instruction (ret-true, ret-zero, etc...)
	literalReturns = globalReplace(
		visitor.Pattern{
			basics.MapKeys(literalReturnMap),
			{isa.Return},
		},
		literalReturnMapper,
	)
)

func literalReturnMapper(i isa.Instructions) isa.Instructions {
	oc := i[0].Opcode()
	res := literalReturnMap[oc]
	return isa.Instructions{res.New()}
}
