package optimize

import (
	"github.com/kode4food/ale/internal/compiler/encoder"
	"github.com/kode4food/ale/internal/compiler/ir/visitor"
	"github.com/kode4food/ale/internal/runtime/isa"
	"github.com/kode4food/ale/internal/runtime/vm"
)

type tailCallMapper struct{ *encoder.Encoded }

var (
	tailCallOpcode = []isa.Opcode{
		isa.Call0, isa.Call1, isa.Call2, isa.Call3, isa.Call,
	}

	tailCallPattern = visitor.Pattern{
		{visitor.AnyOpcode}, tailCallOpcode, {isa.Return},
	}

	selfCallsInTailPosition = globalReplace(
		visitor.Pattern{{isa.CallSelf}, {isa.Return}},
		func(i isa.Instructions) isa.Instructions {
			return isa.Instructions{
				isa.TailSelf.New(getCallArgCount(i[0])),
			}
		},
	)
)

// callsInTailPosition replaces calls in tail position with a tail-call instruction
func callsInTailPosition(e *encoder.Encoded) *encoder.Encoded {
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
		inst = isa.TailClos
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
