package optimize

import (
	"github.com/kode4food/ale/compiler/ir/visitor"
	"github.com/kode4food/ale/runtime/isa"
)

var tailCallPattern = visitor.Pattern{
	{isa.MakeCall},
	{isa.Call, isa.Call0, isa.Call1},
	{isa.Return},
}

func tailCalls(root visitor.Node) visitor.Node {
	visitor.Replace(root, tailCallPattern, tailCallMapper)
	return root
}

func tailCallMapper(i isa.Instructions) isa.Instructions {
	var argCount isa.Word
	switch i[1].Opcode {
	case isa.Call1:
		argCount = 1
	case isa.Call:
		argCount = i[1].Args[0]
	}
	return isa.Instructions{
		isa.New(isa.TailCall, argCount),
	}
}
