package optimize

import (
	"gitlab.com/kode4food/ale/compiler/ir/visitor"
	"gitlab.com/kode4food/ale/runtime/isa"
)

var literalReturnMap = map[isa.Opcode]isa.Opcode{
	isa.True:  isa.RetTrue,
	isa.False: isa.RetFalse,
	isa.Nil:   isa.RetNil,
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
	return isa.Instructions{
		isa.New(literalReturnMap[i[0].Opcode]),
	}
}
