package optimize

import (
	"github.com/kode4food/ale/compiler/ir/visitor"
	"github.com/kode4food/ale/runtime/isa"
)

var unTruthyPattern = visitor.Pattern{
	{isa.True, isa.False, isa.Eq, isa.Neq, isa.Gt, isa.Gte,
		isa.Lt, isa.Lte, isa.Not},
	{isa.MakeTruthy},
}

func unTruthy(root visitor.Node) visitor.Node {
	visitor.Replace(root, unTruthyPattern, unTruthyMapper)
	return root
}

func unTruthyMapper(i isa.Instructions) isa.Instructions {
	return i[0:1]
}
