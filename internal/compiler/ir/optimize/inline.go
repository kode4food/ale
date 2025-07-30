package optimize

import (
	"cmp"
	"slices"

	"github.com/kode4food/ale"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/basics"
	"github.com/kode4food/ale/internal/compiler/encoder"
	"github.com/kode4food/ale/internal/compiler/ir/visitor"
	"github.com/kode4food/ale/internal/runtime/isa"
	"github.com/kode4food/ale/internal/runtime/vm"
)

type (
	inlineMapper struct {
		*encoder.Encoded
		constants data.Vector
		nextLabel isa.Operand
		baseLocal isa.Operand
	}

	operandMap map[isa.Operand]isa.Operand
)

var (
	inlineCallOpcode = []isa.Opcode{
		isa.Call0, isa.Call1, isa.Call2, isa.Call3, isa.Call,
	}

	inlineCallPattern = visitor.Pattern{{isa.Const}, inlineCallOpcode}

	mapReturns = map[isa.Opcode]isa.Instruction{
		isa.RetTrue:  isa.True.New(),
		isa.RetFalse: isa.False.New(),
		isa.RetNull:  isa.Null.New(),
	}
)

// inlineCalls imports callee instructions into procedures that qualify
func inlineCalls(e *encoder.Encoded) *encoder.Encoded {
	m := &inlineMapper{
		Encoded:   e,
		constants: e.Constants,
		nextLabel: getNextLabel(e.Code),
		baseLocal: getNextLocal(e.Code),
	}
	r := visitor.Replace(inlineCallPattern, m.perform)
	res := performReplace(e, r)
	return res.WithConstants(m.constants)
}

func (m *inlineMapper) perform(i isa.Instructions) isa.Instructions {
	p, ok := m.canInline(i[0])
	if !ok {
		return i
	}

	argc := getCallArgCount(i[1])
	c := m.relabel(p.Code)
	c = paramBranchFor(c, argc)
	if hasTailCallInstruction(c) {
		return i
	}
	c = m.reindex(p, c)
	c = m.returns(c)
	c = m.transformArgs(c, argc)
	return c
}

func (m *inlineMapper) canInline(i isa.Instruction) (*vm.Closure, bool) {
	p, ok := m.constants[i.Operand()].(*vm.Closure)
	return p, ok && p.Globals == m.Globals
}

func (m *inlineMapper) relabel(c isa.Instructions) isa.Instructions {
	res := slices.Clone(c)
	labels := operandMap{}
	for idx, i := range res {
		if oc, op := i.Split(); isa.Effects[oc].Operand == isa.Labels {
			to, ok := labels[op]
			if !ok {
				to = m.nextLabel
				m.nextLabel++
				labels[op] = to
			}
			res[idx] = oc.New(to)
		}
	}
	if len(labels) == 0 {
		return res
	}
	s := basics.SortedKeysFunc(labels, func(l, r isa.Operand) int {
		return -cmp.Compare(l, r)
	})
	for _, oc := range s {
		res = slices.Insert(res, int(oc), isa.Label.New(labels[oc]))
	}
	return res
}

func (m *inlineMapper) reindex(
	p *vm.Closure, c isa.Instructions,
) isa.Instructions {
	res := slices.Clone(c)
	pc := p.Captured()
	captured := operandMap{}
	for idx, i := range res {
		switch oc, op := i.Split(); oc {
		case isa.Const:
			val := p.Constants[op]
			to := m.addConstant(val)
			res[idx] = oc.New(to)
		case isa.Closure:
			to, ok := captured[op]
			if !ok {
				to = m.addConstant(pc[op])
				captured[op] = to
			}
			res[idx] = isa.Const.New(to)
		case isa.Load, isa.Store:
			res[idx] = oc.New(op + m.baseLocal)
		default:
			// no-op
		}
	}
	return res
}

func (m *inlineMapper) addConstant(val ale.Value) isa.Operand {
	c := m.constants
	if idx, ok := c.IndexOf(val); ok {
		return isa.Operand(idx)
	}
	c = append(c, val)
	m.constants = c
	return isa.Operand(len(c) - 1)
}

func (m *inlineMapper) returns(c isa.Instructions) isa.Instructions {
	res := make(isa.Instructions, 0, len(c))
	var label isa.Operand
	var increment bool
	if oc, op := c[len(c)-1].Split(); oc == isa.Label {
		label = op
	} else {
		label = m.nextLabel
		increment = true
	}
	for _, i := range c {
		switch oc := i.Opcode(); oc {
		case isa.Return:
			res = append(res, isa.Jump.New(label))
		case isa.RetTrue, isa.RetFalse, isa.RetNull:
			res = append(res, mapReturns[oc], isa.Jump.New(label))
		default:
			res = append(res, i)
		}
	}
	if increment {
		res = append(res, isa.Label.New(label))
		m.nextLabel++
	}
	return res
}

func (m *inlineMapper) transformArgs(
	c isa.Instructions, argc isa.Operand,
) isa.Instructions {
	switch argc {
	case 0:
		if !hasAnyArgInstruction(c) {
			return c
		}
	case 1:
		if res, ok := immediatelyPopArg(c); ok {
			return res
		}
		fallthrough
	default:
		if canMapArgsToLocals(c, argc) {
			argsBase := m.baseLocal + getNextLocal(c)
			return mapArgsToLocals(c, argsBase, argc)
		}
	}
	return stackArgs(c, argc)
}

func getNextLabel(c isa.Instructions) isa.Operand {
	return getNextOperand(c, isa.Labels)
}

func getNextLocal(c isa.Instructions) isa.Operand {
	return getNextOperand(c, isa.Locals)
}

func getNextOperand(c isa.Instructions, actOn isa.ActOn) isa.Operand {
	var res isa.Operand
	for _, i := range c {
		oc, op := i.Split()
		if isa.Effects[oc].Operand == actOn && op >= res {
			res = op + 1
		}
	}
	return res
}

func immediatelyPopArg(c isa.Instructions) (isa.Instructions, bool) {
	if len(c) != 0 || c[0] != isa.Arg.New(0) {
		return nil, false
	}
	if res := c[1:]; !hasAnyArgInstruction(res) {
		return res, true
	}
	return nil, false
}

func hasTailCallInstruction(c isa.Instructions) bool {
	return slices.ContainsFunc(c, func(i isa.Instruction) bool {
		switch i.Opcode() {
		case isa.TailCall, isa.TailClos:
			return true
		default:
			return false
		}
	})
}
