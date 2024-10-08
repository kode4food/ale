package optimize

import (
	"github.com/kode4food/ale/internal/compiler/encoder"
	"github.com/kode4food/ale/internal/compiler/ir/visitor"
	"github.com/kode4food/ale/internal/debug"
	"github.com/kode4food/ale/internal/runtime/isa"
	"github.com/kode4food/ale/internal/runtime/vm"
)

type tailCallMapper struct{ *encoder.Encoded }

var tailCallPattern = visitor.Pattern{
	{visitor.AnyOpcode},
	{isa.Call, isa.Call0, isa.Call1},
	{isa.Return},
}

// makeTailCalls replaces calls in tail position with a tail-call instruction
func makeTailCalls(e *encoder.Encoded) *encoder.Encoded {
	m := &tailCallMapper{e}
	r := visitor.Replace(tailCallPattern, m.perform)
	return performReplace(e, r)
}

func (m tailCallMapper) perform(i isa.Instructions) isa.Instructions {
	if !m.canTailCall(i[0]) {
		return i
	}
	var argc isa.Operand
	switch oc, op := i[1].Split(); oc {
	case isa.Call0:
		// no-op
	case isa.Call1:
		argc = 1
	case isa.Call:
		argc = op
	default:
		panic(debug.ProgrammerError("bad opcode matching"))
	}
	return isa.Instructions{
		i[0],
		isa.TailCall.New(argc),
	}
}

func (m tailCallMapper) canTailCall(i isa.Instruction) bool {
	if oc, op := i.Split(); oc == isa.Const {
		_, ok := m.Constants[op].(*vm.Closure)
		return ok
	}
	return true
}
