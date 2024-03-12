package encoder

import (
	"errors"
	"slices"

	"github.com/kode4food/ale/compiler/ir/analysis"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/env"
	"github.com/kode4food/ale/runtime/isa"
)

type (
	// Encoded is a snapshot of the current Encoder's state. It is used as an
	// intermediate step in the compilation process, particularly as input to
	// the optimizer.
	Encoded struct {
		Code      isa.Instructions
		Globals   env.Namespace
		Constants data.Vector
		Closure   data.Locals
	}

	finalizer struct {
		*Encoded
		labels     map[isa.Operand]*label
		constMap   map[isa.Operand]isa.Operand
		output     isa.Instructions
		constants  data.Vector
		localCount isa.Operand
	}

	label struct {
		backRefs []isa.Operand
		offset   isa.Operand
		anchored bool
	}
)

// ErrLabelAlreadyAnchored is raised when the finalizer identifies that a label
// has been anchored more than once in the Instructions stream
const ErrLabelAlreadyAnchored = "label has already been anchored"

func (e *Encoded) Copy() *Encoded {
	res := *e
	res.Code = slices.Clone(e.Code)
	res.Constants = slices.Clone(e.Constants)
	res.Closure = slices.Clone(e.Closure)
	return &res
}

// Runnable takes an Encoded and finalizes it into a Runnable that the abstract
// machine can execute. Jumps are resolved and unused constants are discarded.
func (e *Encoded) Runnable() *isa.Runnable {
	f := &finalizer{
		Encoded:  e,
		labels:   map[isa.Operand]*label{},
		constMap: map[isa.Operand]isa.Operand{},
	}
	return f.finalize()
}

func (f *finalizer) finalize() *isa.Runnable {
	for _, inst := range f.Code {
		f.handleInst(inst)
	}
	f.stripAdjacentJumps()
	stackSize := analysis.MustCalculateStackSize(f.output)
	return &isa.Runnable{
		Code:       f.output,
		Globals:    f.Globals,
		Constants:  f.constants,
		LocalCount: f.localCount,
		StackSize:  stackSize,
	}
}

func (f *finalizer) stripAdjacentJumps() {
	for i := len(f.output) - 1; i >= 0; {
		oc, op := f.output[i].Split()
		if oc == isa.Jump && op == isa.Operand(i+1) {
			f.output = removeInstruction(f.output, i)
			continue
		}
		i--
	}
}

func removeInstruction(inst isa.Instructions, idx int) isa.Instructions {
	res := append(inst[:idx], inst[idx+1:]...)
	for j, inst := range res {
		oc, op := inst.Split()
		if (oc == isa.Jump || oc == isa.CondJump) && op > isa.Operand(idx) {
			res[j] = isa.New(oc, op-1)
		}
	}
	return res
}

func (f *finalizer) handleInst(i isa.Instruction) {
	switch oc := i.Opcode(); oc {
	case isa.Load, isa.Store:
		f.handleLocal(i)
	case isa.Const:
		f.handleConst(i)
	case isa.Jump:
		f.handleJump(i)
	case isa.CondJump:
		f.handleJump(i)
	case isa.Label:
		f.handleLabel(i)
	default:
		if effect := isa.MustGetEffect(oc); effect.Ignore {
			return
		}
		f.output = append(f.output, i)
	}
}

func (f *finalizer) handleLocal(i isa.Instruction) {
	if op := i.Operand(); op >= f.localCount {
		f.localCount = op + 1
	}
	f.output = append(f.output, i)
}

func (f *finalizer) handleConst(i isa.Instruction) {
	op := i.Operand()
	if idx, ok := f.constMap[op]; ok {
		ni := isa.New(isa.Const, idx)
		f.output = append(f.output, ni)
		return
	}
	idx := isa.Operand(len(f.constants))
	f.constants = append(f.constants, f.Constants[op])
	f.constMap[op] = idx
	ni := isa.New(isa.Const, idx)
	f.output = append(f.output, ni)
}

func (f *finalizer) handleJump(i isa.Instruction) {
	oc, op := i.Split()
	lbl := f.getLabel(op)
	if !lbl.anchored {
		f.addLabelBackRef(lbl)
	}
	ni := isa.New(oc, lbl.offset)
	f.output = append(f.output, ni)
}

func (f *finalizer) handleLabel(i isa.Instruction) {
	op := i.Operand()
	lbl := f.getLabel(op)
	if lbl.anchored {
		panic(errors.New(ErrLabelAlreadyAnchored))
	}
	lbl.offset = f.nextOutputOffset()
	lbl.anchored = true
	for _, off := range lbl.backRefs {
		oc := f.output[int(off)].Opcode()
		ni := isa.New(oc, lbl.offset)
		f.output[int(off)] = ni
	}
}

func (f *finalizer) getLabel(idx isa.Operand) *label {
	if lbl, ok := f.labels[idx]; ok {
		return lbl
	}
	lbl := new(label)
	f.labels[idx] = lbl
	return lbl
}

func (f *finalizer) nextOutputOffset() isa.Operand {
	return isa.Operand(len(f.output))
}

func (f *finalizer) addLabelBackRef(l *label) {
	off := f.nextOutputOffset()
	l.backRefs = append(l.backRefs, off)
}
