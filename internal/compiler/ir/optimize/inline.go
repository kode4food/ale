package optimize

import (
	"cmp"
	"slices"

	"github.com/kode4food/ale/internal/compiler/encoder"
	"github.com/kode4food/ale/internal/compiler/ir/visitor"
	"github.com/kode4food/ale/internal/debug"
	"github.com/kode4food/ale/internal/runtime/isa"
	"github.com/kode4food/ale/internal/runtime/vm"
	"github.com/kode4food/ale/pkg/data"
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
	replace := visitor.Replace(inlineCallPattern, mapper.perform)
	visitor.Visit(root, replace)
	e.Code = root.Code()
}

func (m *inlineCallMapper) perform(i isa.Instructions) isa.Instructions {
	p, ok := m.canInline(i[0])
	if !ok {
		return i
	}
	m.numInlined++

	argCount := getCallArgCount(i[1])

	c := m.relabel(p.Code)
	c = m.getParamCase(c, argCount)
	c = m.reindex(p, c)
	c = m.returns(c)

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

func (m *inlineCallMapper) getParamCase(
	c isa.Instructions, argCount isa.Operand,
) isa.Instructions {
	b := &visitor.BranchScanner{
		Then:     visitor.All,
		Epilogue: visitor.All,
	}
	b.Else = b.Scan

	if b, ok := b.Scan(c).(visitor.Branches); ok {
		if bc := getParamBranch(b, argCount); bc != nil {
			return bc
		}
	}
	return c
}

func (m *inlineCallMapper) reindex(
	p *vm.Closure, c isa.Instructions,
) isa.Instructions {
	res := slices.Clone(c)
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
		default:
			// No-Op
		}
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

func getParamBranch(b visitor.Branches, argCount isa.Operand) isa.Instructions {
	oc, op, ok := isParamCase(b)
	if !ok {
		return nil
	}
	switch {
	case oc == isa.NumEq && argCount == op:
		return b.ThenBranch().Code()
	case oc == isa.NumGte && argCount >= op:
		return b.ThenBranch().Code()
	default:
		if eb, ok := b.ElseBranch().(visitor.Branches); ok {
			return getParamBranch(eb, argCount)
		}
		return nil
	}
}

func isParamCase(b visitor.Branches) (isa.Opcode, isa.Operand, bool) {
	p := b.Prologue().Code()
	if len(p) != 4 {
		return isa.NoOp, 0, false
	}
	if p[0].Opcode() != isa.ArgLen || p[3].Opcode() != isa.CondJump {
		return isa.NoOp, 0, false
	}
	if p[1].Opcode() != isa.PosInt {
		return isa.NoOp, 0, false
	}
	oc := p[2].Opcode()
	op := p[1].Operand()
	return oc, op, oc == isa.NumEq || oc == isa.NumGte
}
