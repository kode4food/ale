package isa

import "errors"

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

// Error messages
const (
	ErrLabelAlreadyAnchored = "label has already been anchored"
)

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
	res := make(Instructions, len(f.output))
	copy(res, f.output)
	return res
}

func (f *flattener) handleInst(l Instruction) {
	oc, _ := l.Split()
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
	_, op := inst.Split()
	l := f.getLabel(op)
	if l.anchored {
		panic(errors.New(ErrLabelAlreadyAnchored))
	}
	l.offset = f.nextOutputOffset()
	l.anchored = true
	for _, off := range l.backRefs {
		oc, _ := f.output[int(off)].Split()
		ni := New(oc, l.offset)
		f.output[int(off)] = ni
	}
}
