package optimize

import (
	"github.com/kode4food/ale/compiler/encoder"
	"github.com/kode4food/ale/compiler/ir/visitor"
	"github.com/kode4food/ale/runtime/isa"
	"github.com/kode4food/comb/basics"
)

var (
	literalReturnMap = map[isa.Opcode]isa.Opcode{
		isa.False: isa.RetFalse,
		isa.Null:  isa.RetNull,
		isa.True:  isa.RetTrue,
	}

	literalReturnPatterns = visitor.Pattern{
		basics.MapKeys(literalReturnMap),
		{isa.Return},
	}
)

func makeLiteralReturns(*encoder.Encoded) optimizer {
	return func(root visitor.Node) visitor.Node {
		visitor.Replace(root, literalReturnPatterns, literalReturnMapper)
		return root
	}
}

func literalReturnMapper(i isa.Instructions) isa.Instructions {
	oc := i[0].Opcode()
	res := literalReturnMap[oc]
	return isa.Instructions{res.New()}
}
