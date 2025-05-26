package optimize

import (
	"github.com/kode4food/ale/internal/compiler/encoder"
	"github.com/kode4food/ale/internal/compiler/ir/visitor"
	"github.com/kode4food/ale/internal/runtime/isa"
	"github.com/kode4food/ale/internal/runtime/vm"
)

type tailCallMapper struct{ *encoder.Encoded }

var tailCallPattern = visitor.Pattern{
	{visitor.AnyOpcode}, anyCallOpcode, {isa.Return},
}

// makeTailCalls replaces calls in tail position with a tail-call instruction
func makeTailCalls(e *encoder.Encoded) *encoder.Encoded {
	m := &tailCallMapper{e}
	r := visitor.Replace(tailCallPattern, m.perform)
	return performReplace(e, r)
}

func (m tailCallMapper) perform(i isa.Instructions) isa.Instructions {
	c, ok := m.canTailCall(i[0])
	if !ok {
		return i
	}
	inst := isa.TailCall
	if c != nil {
		inst = isa.TailDiff
	}
	return isa.Instructions{
		i[0],
		inst.New(getCallArgCount(i[1])),
	}
}

func (m tailCallMapper) canTailCall(i isa.Instruction) (*vm.Closure, bool) {
	if oc, op := i.Split(); oc == isa.Const {
		c, ok := m.Constants[op].(*vm.Closure)
		return c, ok
	}
	return nil, true
}
