package isa

import (
	"errors"
	"slices"
)

type (
	flattener struct {
		labels labels
		input  Instructions
		output Instructions
	}

	label struct {
		anchored bool
		offset   Operand
		backRefs []Operand
	}

	labels map[Operand]*label
)

// ErrLabelAlreadyAnchored is raised when the flattener identifies that a label
// has been anchored more than once in the Instructions stream
const ErrLabelAlreadyAnchored = "label has already been anchored"

// Flatten takes a set of Instructions and flattens them into something that
// the abstract machine can execute
func Flatten(code Instructions) Instructions {
	f := &flattener{
		input:  code,
		labels: labels{},
	}
	return f.flatten()
}

func (f *flattener) flatten() Instructions {
	for _, l := range f.input {
		f.handleInst(l)
	}
	return slices.Clone(f.output)
}

func (f *flattener) handleInst(l Instruction) {
	oc := l.Opcode()
	switch oc {
	case Jump:
		f.handleJump(l)
	case CondJump:
		f.handleJump(l)
	case Label:
		f.handleLabel(l)
	default:
		if effect := MustGetEffect(oc); effect.Ignore {
			return
		}
		f.output = append(f.output, l)
	}
}

func (f *flattener) getLabel(idx Operand) *label {
	if l, ok := f.labels[idx]; ok {
		return l
	}
	l := new(label)
	f.labels[idx] = l
	return l
}

func (f *flattener) nextOutputOffset() Operand {
	return Operand(len(f.output))
}

func (f *flattener) addLabelBackRef(l *label) {
	off := f.nextOutputOffset()
	l.backRefs = append(l.backRefs, off)
}

func (f *flattener) handleJump(inst Instruction) {
	oc, op := inst.Split()
	l := f.getLabel(op)
	if !l.anchored {
		f.addLabelBackRef(l)
	}
	ni := New(oc, l.offset)
	f.output = append(f.output, ni)
}

func (f *flattener) handleLabel(inst Instruction) {
	op := inst.Operand()
	l := f.getLabel(op)
	if l.anchored {
		panic(errors.New(ErrLabelAlreadyAnchored))
	}
	l.offset = f.nextOutputOffset()
	l.anchored = true
	for _, off := range l.backRefs {
		oc := f.output[int(off)].Opcode()
		ni := New(oc, l.offset)
		f.output[int(off)] = ni
	}
}
