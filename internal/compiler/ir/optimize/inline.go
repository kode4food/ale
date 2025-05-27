package optimize

import (
	"cmp"
	"slices"

	"github.com/kode4food/ale/internal/basics"
	"github.com/kode4food/ale/internal/compiler/encoder"
	"github.com/kode4food/ale/internal/compiler/ir/visitor"
	"github.com/kode4food/ale/internal/debug"
	"github.com/kode4food/ale/internal/runtime/isa"
	"github.com/kode4food/ale/internal/runtime/vm"
	"github.com/kode4food/ale/pkg/data"
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

func (m *inlineMapper) addConstant(val data.Value) isa.Operand {
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
		if res, ok := immediatelyPop(c); ok {
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

func getCallArgCount(i isa.Instruction) isa.Operand {
	switch i.Opcode() {
	case isa.Call0:
		return 0
	case isa.Call1:
		return 1
	case isa.Call2:
		return 2
	case isa.Call3:
		return 3
	case isa.Call, isa.CallSelf:
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

func paramBranchFor(c isa.Instructions, argc isa.Operand) isa.Instructions {
	b := &visitor.BranchScanner{
		Then:     visitor.All,
		Epilogue: visitor.All,
	}
	b.Else = b.Scan

	if b, ok := b.Scan(c).(visitor.Branches); ok {
		if bc := getParamBranch(b, argc); bc != nil {
			return bc
		}
	}
	return c
}

func getParamBranch(b visitor.Branches, argc isa.Operand) isa.Instructions {
	if len(b.Epilogue().Code()) != 0 {
		// compiled procedures don't include epilogues in the arity branching
		// logic, so if any node along the path has an epilogue, then we can't
		// inline the 'then' branch
		return nil
	}
	oc, op, ok := isParamCase(b)
	if !ok {
		return nil
	}
	if oc == isa.NumEq && argc == op || oc == isa.NumGte && argc >= op {
		return b.ThenBranch().Code()
	}
	if eb, ok := b.ElseBranch().(visitor.Branches); ok {
		return getParamBranch(eb, argc)
	}
	return nil
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

func hasTailCallInstruction(c isa.Instructions) bool {
	return slices.ContainsFunc(c, func(i isa.Instruction) bool {
		switch i.Opcode() {
		case isa.TailCall, isa.TailDiff:
			return true
		default:
			return false
		}
	})
}

func hasAnyArgInstruction(c isa.Instructions) bool {
	return slices.ContainsFunc(c, argInstructionPred)
}

func filterArgInstructions(c isa.Instructions) isa.Instructions {
	return basics.Filter(c, argInstructionPred)
}

func argInstructionPred(i isa.Instruction) bool {
	switch i.Opcode() {
	case isa.PushArgs, isa.PopArgs, isa.Arg, isa.ArgLen, isa.RestArg:
		return true
	default:
		return false
	}
}

func immediatelyPop(c isa.Instructions) (isa.Instructions, bool) {
	if len(c) != 0 || c[0] != isa.Arg.New(0) {
		return nil, false
	}
	if res := c[1:]; !hasAnyArgInstruction(res) {
		return res, true
	}
	return nil, false
}

func canMapArgsToLocals(c isa.Instructions, argc isa.Operand) bool {
	highArg := -1
	for _, i := range filterArgInstructions(c) {
		switch i.Opcode() {
		case isa.Arg:
			idx := int(i.Operand())
			if idx > highArg {
				highArg = idx
			}
		case isa.PushArgs, isa.PopArgs, isa.ArgLen, isa.RestArg:
			return false
		default:
			// no-op
		}
	}
	return highArg < int(argc)
}

func mapArgsToLocals(
	c isa.Instructions, argsBase, argc isa.Operand,
) isa.Instructions {
	al := makeArgLocalMap(c, argsBase)
	res := make(isa.Instructions, 0, len(c)+int(argc))
	for i := range int(argc) {
		if to, ok := al[isa.Operand(i)]; ok {
			res = append(res, isa.Store.New(to))
			continue
		}
		res = append(res, isa.Pop.New())
	}
	res = append(res, basics.Map(c, func(i isa.Instruction) isa.Instruction {
		if i.Opcode() == isa.Arg {
			to := al[i.Operand()]
			return isa.Load.New(to)
		}
		return i
	})...)
	return res
}

func makeArgLocalMap(c isa.Instructions, argsBase isa.Operand) operandMap {
	next := argsBase
	res := operandMap{}
	for _, i := range c {
		if i.Opcode() != isa.Arg {
			continue
		}
		op := i.Operand()
		if _, ok := res[op]; !ok {
			res[op] = next
			next++
		}
	}
	return res
}

func stackArgs(c isa.Instructions, argc isa.Operand) isa.Instructions {
	res := make(isa.Instructions, 0, len(c)+2)
	res = append(res, isa.PushArgs.New(argc))
	res = append(res, c...)
	res = append(res, isa.PopArgs.New())
	return res
}
