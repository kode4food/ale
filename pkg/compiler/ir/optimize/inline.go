package optimize

import (
	"cmp"
	"slices"

	"github.com/kode4food/ale/internal/debug"
	"github.com/kode4food/ale/pkg/compiler/encoder"
	"github.com/kode4food/ale/pkg/compiler/ir/visitor"
	"github.com/kode4food/ale/pkg/data"
	"github.com/kode4food/ale/pkg/runtime/isa"
	"github.com/kode4food/ale/pkg/runtime/vm"
	"github.com/kode4food/comb/basics"
)

type inlineCallMapper struct {
	*encoder.Encoded
	numInlined int
	nextLabel  isa.Operand
	baseLocal  isa.Operand
}

const maxInlined = 16

var (
	inlineCallPattern = visitor.Pattern{
		{isa.Const},
		{isa.Call0, isa.Call1, isa.Call},
	}

	mapReturns = map[isa.Opcode]isa.Instruction{
		isa.RetTrue:  isa.True.New(),
		isa.RetFalse: isa.False.New(),
		isa.RetNull:  isa.Null.New(),
	}
)

func inlineCalls(e *encoder.Encoded) {
	mapper := &inlineCallMapper{
		Encoded:   e,
		nextLabel: getNextLabel(e.Code),
		baseLocal: getNextLocal(e.Code),
	}
	root := visitor.All(e.Code)
	visitor.Replace(inlineCallPattern, mapper.perform).Instructions(root)
	e.Code = root.Code()
}

func (m *inlineCallMapper) perform(i isa.Instructions) isa.Instructions {
	p, ok := m.canInline(i[0])
	if !ok {
		return i
	}
	m.numInlined++

	c := m.reindex(p)
	c = m.relabel(c)
	c = m.returns(c)

	argCount := getCallArgCount(i[1])
	argsLocal := m.baseLocal + getNextLocal(p.Code)
	c = m.stackArgs(c, argCount, argsLocal)

	return c
}

func (m *inlineCallMapper) canInline(i isa.Instruction) (*vm.Closure, bool) {
	p, ok := m.Constants[i.Operand()].(*vm.Closure)
	return p, ok &&
		!p.HasFlag(vm.NoInline) &&
		m.numInlined < maxInlined &&
		p.Globals == m.Globals
}

func (m *inlineCallMapper) reindex(p *vm.Closure) isa.Instructions {
	res := slices.Clone(p.Code)
	pc := p.Captured()
	captured := map[isa.Operand]isa.Operand{}
	for idx, i := range res {
		switch oc, op := i.Split(); oc {
		case isa.Const:
			val := p.Constants[op]
			to := m.addConstant(val)
			res[idx] = isa.New(oc, to)
		case isa.Closure:
			to, ok := captured[op]
			if !ok {
				to = m.addConstant(pc[op])
				captured[op] = to
			}
			res[idx] = isa.New(isa.Const, to)
		case isa.Load, isa.Store:
			res[idx] = isa.New(oc, op+m.baseLocal)
		}
	}
	return res
}

func (m *inlineCallMapper) relabel(c isa.Instructions) isa.Instructions {
	res := slices.Clone(c)
	labels := map[isa.Operand]isa.Operand{}
	for idx, i := range res {
		if oc, op := i.Split(); isa.Effects[oc].Operand == isa.Labels {
			to, ok := labels[op]
			if !ok {
				to = m.nextLabel
				m.nextLabel++
				labels[op] = to
			}
			res[idx] = isa.New(oc, to)
		}
	}
	if len(labels) == 0 {
		return res
	}
	s := basics.SortedKeysFunc(labels, func(l, r isa.Operand) int {
		return -cmp.Compare(l, r)
	})
	for _, oc := range s {
		res = slices.Insert(res, int(oc), isa.New(isa.Label, labels[oc]))
	}
	return res
}

func (m *inlineCallMapper) returns(c isa.Instructions) isa.Instructions {
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

func (m *inlineCallMapper) addConstant(val data.Value) isa.Operand {
	c := m.Constants
	if idx, ok := c.IndexOf(val); ok {
		return isa.Operand(idx)
	}
	c = append(c, val)
	m.Constants = c
	return isa.Operand(len(c) - 1)
}

func (m *inlineCallMapper) stackArgs(
	c isa.Instructions, argc isa.Operand, argsLocal isa.Operand,
) isa.Instructions {
	res := make(isa.Instructions, 0, len(c)+6)
	res = append(res,
		isa.RestArg.New(0),
		isa.Store.New(argsLocal),
		isa.Vector.New(argc),
		isa.SetArgs.New(),
	)
	res = append(res, c...)
	res = append(res,
		isa.Load.New(argsLocal),
		isa.SetArgs.New(),
	)
	return res
}

func getCallArgCount(i isa.Instruction) isa.Operand {
	switch i.Opcode() {
	case isa.Call0:
		return 0
	case isa.Call1:
		return 1
	case isa.Call:
		return i.Operand()
	default:
		panic(debug.ProgrammerError("invalid call instruction matched"))
	}
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
