package optimize

import (
	"github.com/kode4food/ale/compiler/ir/visitor"
	"github.com/kode4food/ale/runtime/isa"
)

var stripTruthyPattern = visitor.Pattern{
	{isa.True, isa.False, isa.NumEq, isa.NumNeq, isa.NumGt, isa.NumGte,
		isa.NumLt, isa.NumLte, isa.Not},
	{isa.MakeTruthy},
}

func stripTruthy(root visitor.Node) visitor.Node {
	visitor.Replace(root, stripTruthyPattern, stripTruthyMapper)
	return root
}

func stripTruthyMapper(i isa.Instructions) isa.Instructions {
	return i[0:1]
}
