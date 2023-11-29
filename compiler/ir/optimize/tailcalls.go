package optimize

import (
	"github.com/kode4food/ale/compiler/encoder"
	"github.com/kode4food/ale/compiler/ir/visitor"
	"github.com/kode4food/ale/runtime/isa"
	"github.com/kode4food/ale/runtime/vm"
)

type tailCallMapper struct{ encoder.Encoder }

var tailCallPattern = visitor.Pattern{
	{visitor.AnyOpcode},
	{isa.Call, isa.Call0, isa.Call1},
	{isa.Return},
}

func makeTailCalls(e encoder.Encoder) optimizer {
	return func(root visitor.Node) visitor.Node {
		visitor.Replace(root, tailCallPattern, tailCallMapper{e}.perform)
		return root
	}
}

func (m tailCallMapper) canTailCall(i isa.Instruction) bool {
	if oc, op := i.Split(); oc == isa.Const {
		_, ok := m.Constants()[op].(*vm.Closure)
		return ok
	}
	return true
}

func (m tailCallMapper) perform(i isa.Instructions) isa.Instructions {
	if !m.canTailCall(i[0]) {
		return i
	}
	var argCount isa.Operand
	oc, op := i[1].Split()
	switch oc {
	case isa.Call1:
		argCount = 1
	case isa.Call:
		argCount = op
	}
	return isa.Instructions{
		i[0],
		isa.TailCall.New(argCount),
	}
}
