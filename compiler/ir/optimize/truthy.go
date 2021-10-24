package optimize

import (
	"github.com/kode4food/ale/compiler/ir/visitor"
	"github.com/kode4food/ale/runtime/isa"
)

var stripTruthyPattern = visitor.Pattern{
	{isa.True, isa.False, isa.Eq, isa.Neq, isa.Gt, isa.Gte,
		isa.Lt, isa.Lte, isa.Not},
	{isa.MakeTruthy},
}

func stripTruthy(root visitor.Node) visitor.Node {
	visitor.Replace(root, stripTruthyPattern, stripTruthyMapper)
	return root
}

func stripTruthyMapper(i isa.Instructions) isa.Instructions {
	return i[0:1]
}
