package optimize

import (
	"gitlab.com/kode4food/ale/compiler/ir/visitor"
	"gitlab.com/kode4food/ale/runtime/isa"
)

var literalReturnMap = map[isa.Opcode]isa.Opcode{
	isa.False: isa.RetFalse,
	isa.Null:  isa.RetNull,
	isa.True:  isa.RetTrue,
}

var literalReturnPatterns = [][]isa.Opcode{
	literalKeys(),
	{isa.Return},
}

func literalKeys() []isa.Opcode {
	var res []isa.Opcode
	for k := range literalReturnMap {
		res = append(res, k)
	}
	return res
}

func literalReturns(root visitor.Node) visitor.Node {
	visitor.Replace(root, literalReturnPatterns, literalReturnMapper)
	return root
}

func literalReturnMapper(i isa.Instructions) isa.Instructions {
	oc := i[0].Opcode
	res := literalReturnMap[oc]
	return isa.Instructions{
		isa.New(res),
	}
}
