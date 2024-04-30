package optimize

import (
	"github.com/kode4food/ale/internal/compiler/encoder"
	"github.com/kode4food/ale/internal/compiler/ir/visitor"
	"github.com/kode4food/ale/internal/runtime/isa"
	"github.com/kode4food/comb/basics"
)

var (
	literalReturnMap = map[isa.Opcode]isa.Opcode{
		isa.False: isa.RetFalse,
		isa.Null:  isa.RetNull,
		isa.True:  isa.RetTrue,
	}

	literalReturnReplace = visitor.Replace(
		visitor.Pattern{
			basics.MapKeys(literalReturnMap),
			{isa.Return},
		},
		literalReturnMapper,
	)
)

func literalReturns(e *encoder.Encoded) {
	root := visitor.All(e.Code)
	literalReturnReplace.Instructions(root)
	e.Code = root.Code()
}

func literalReturnMapper(i isa.Instructions) isa.Instructions {
	oc := i[0].Opcode()
	res := literalReturnMap[oc]
	return isa.Instructions{res.New()}
}
