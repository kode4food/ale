package isa

import "math"

type (
	flattener struct {
		labels labels
		input  Instructions
		output []Word
	}

	label struct {
		anchored bool
		offset   Offset
		backRefs []Word
	}

	labels map[Index]*label
)

// Error messages
const (
	ErrLabelAlreadyAnchored = "label has already been anchored"
)

const placeholderOffset = Offset(math.MaxUint32)

// Flatten takes a set of instructions and flattens them into
// something that the virtual machine can execute
func Flatten(code Instructions) []Word {
	f := &flattener{
		input:  code,
		output: []Word{},
		labels: labels{},
	}
	return f.flatten()
}

func (f *flattener) flatten() []Word {
	for _, l := range f.input {
		f.handleInst(l)
	}
	res := make([]Word, len(f.output))
	copy(res, f.output)
	return res
}

func (f *flattener) handleInst(l *Instruction) {
	oc := l.Opcode
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
		f.output = append(f.output, Word(oc))
		f.output = append(f.output, l.Args...)
	}
}

func (f *flattener) getLabel(idx Index) *label {
	if l, ok := f.labels[idx]; ok {
		return l
	}
	l := &label{
		offset:   placeholderOffset,
		backRefs: []Word{},
	}
	f.labels[idx] = l
	return l
}

func (f *flattener) nextOutputOffset() Offset {
	return Offset(len(f.output))
}

func (f *flattener) addLabelBackRef(l *label) {
	off := f.nextOutputOffset()
	l.backRefs = append(l.backRefs, Word(off))
}

func (f *flattener) handleJump(inst *Instruction) {
	l := f.getLabel(Index(inst.Args[0]))
	f.output = append(f.output, Word(inst.Opcode))
	if !l.anchored {
		f.addLabelBackRef(l)
	}
	f.output = append(f.output, Word(l.offset))
}

func (f *flattener) handleLabel(inst *Instruction) {
	l := f.getLabel(Index(inst.Args[0]))
	if l.anchored {
		panic(ErrLabelAlreadyAnchored)
	}
	l.offset = f.nextOutputOffset()
	l.anchored = true
	backRefs := l.backRefs
	if len(backRefs) > 0 {
		for _, off := range backRefs {
			f.output[int(off)] = Word(l.offset)
		}
		l.backRefs = nil
	}
}
