package optimize

import (
	"github.com/kode4food/ale/compiler/ir/visitor"
	"github.com/kode4food/ale/internal/maps"
	"github.com/kode4food/ale/runtime/isa"
)

var (
	literalReturnMap = map[isa.Opcode]isa.Opcode{
		isa.False: isa.RetFalse,
		isa.Null:  isa.RetNull,
		isa.True:  isa.RetTrue,
	}

	literalKeys = maps.Keys(literalReturnMap)

	literalReturnPatterns = visitor.Pattern{
		literalKeys,
		{isa.Return},
	}
)

func literalReturns(root visitor.Node) visitor.Node {
	visitor.Replace(root, literalReturnPatterns, literalReturnMapper)
	return root
}

func literalReturnMapper(i isa.Instructions) isa.Instructions {
	oc, _ := i[0].Split()
	res := literalReturnMap[oc]
	return isa.Instructions{res.New()}
}
