package optimize

import (
	"github.com/kode4food/ale/compiler/ir/visitor"
	"github.com/kode4food/ale/runtime/isa"
)

var (
	literalReturnMap = map[isa.Opcode]isa.Opcode{
		isa.False: isa.RetFalse,
		isa.Null:  isa.RetNil,
		isa.True:  isa.RetTrue,
	}

	literalKeys = _makeLiteralKeys()

	literalReturnPatterns = visitor.Pattern{
		literalKeys,
		{isa.Return},
	}
)

func _makeLiteralKeys() []isa.Opcode {
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
	oc, _ := i[0].Split()
	res := literalReturnMap[oc]
	return isa.Instructions{res.New()}
}
