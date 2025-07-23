package encoder

import (
	"errors"
	"slices"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/env"
	"github.com/kode4food/ale/internal/compiler/ir/analysis"
	"github.com/kode4food/ale/internal/runtime/isa"
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
		labels    map[isa.Operand]*label
		constMap  map[isa.Operand]isa.Operand
		localMap  map[isa.Operand]isa.Operand
		output    isa.Instructions
		constants data.Vector
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

func (e *Encoded) WithCode(c isa.Instructions) *Encoded {
	res := *e
	res.Code = c
	return &res
}

func (e *Encoded) WithConstants(c data.Vector) *Encoded {
	res := *e
	res.Constants = c
	return &res
}

func (e *Encoded) HasClosure() bool {
	return len(e.Closure) > 0
}

// Runnable takes an Encoded and finalizes it into a Runnable that the abstract
// machine can execute. Jumps are resolved and unused constants are discarded.
func (e *Encoded) Runnable() (*isa.Runnable, error) {
	f := &finalizer{
		Encoded:  e,
		labels:   map[isa.Operand]*label{},
		constMap: map[isa.Operand]isa.Operand{},
		localMap: map[isa.Operand]isa.Operand{},
	}
	return f.finalize()
}

func (f *finalizer) finalize() (*isa.Runnable, error) {
	for _, inst := range f.Code {
		if err := f.handleInst(inst); err != nil {
			return nil, err
		}
	}
	f.stripAdjacentJumps()
	return f.makeRunnable()
}

func (f *finalizer) handleInst(i isa.Instruction) error {
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
		return f.handleLabel(i)
	default:
		effect, err := isa.GetEffect(oc)
		if err != nil {
			return err
		}
		if !effect.Ignore {
			f.output = append(f.output, i)
		}
	}
	return nil
}

func (f *finalizer) handleLocal(i isa.Instruction) {
	from := i.Operand()
	to, ok := f.localMap[from]
	if !ok {
		to = isa.Operand(len(f.localMap))
		f.localMap[from] = to
	}
	f.output = append(f.output, i.Opcode().New(to))
}

func (f *finalizer) handleConst(i isa.Instruction) {
	op := i.Operand()
	if idx, ok := f.constMap[op]; ok {
		ni := isa.Const.New(idx)
		f.output = append(f.output, ni)
		return
	}
	idx := isa.Operand(len(f.constants))
	f.constants = append(f.constants, f.Constants[op])
	f.constMap[op] = idx
	ni := isa.Const.New(idx)
	f.output = append(f.output, ni)
}

func (f *finalizer) handleJump(i isa.Instruction) {
	oc, op := i.Split()
	lbl := f.getLabel(op)
	if !lbl.anchored {
		f.addLabelBackRef(lbl)
	}
	ni := oc.New(lbl.offset)
	f.output = append(f.output, ni)
}

func (f *finalizer) handleLabel(i isa.Instruction) error {
	op := i.Operand()
	lbl := f.getLabel(op)
	if lbl.anchored {
		return errors.New(ErrLabelAlreadyAnchored)
	}
	lbl.offset = f.nextOutputOffset()
	lbl.anchored = true
	for _, off := range lbl.backRefs {
		oc := f.output[int(off)].Opcode()
		ni := oc.New(lbl.offset)
		f.output[int(off)] = ni
	}
	return nil
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
	res := slices.Concat(inst[:idx], inst[idx+1:])
	for j, inst := range res {
		oc, op := inst.Split()
		if (oc == isa.Jump || oc == isa.CondJump) && op > isa.Operand(idx) {
			res[j] = oc.New(op - 1)
		}
	}
	return res
}

func (f *finalizer) makeRunnable() (*isa.Runnable, error) {
	stackSize, err := analysis.CalculateStackSize(f.output)
	if err != nil {
		return nil, err
	}
	return &isa.Runnable{
		Code:       f.output,
		Globals:    f.Globals,
		Constants:  f.constants,
		LocalCount: isa.Operand(len(f.localMap)),
		StackSize:  stackSize,
	}, nil
}
