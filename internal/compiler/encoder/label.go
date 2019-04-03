package encoder

import (
	"math"

	"gitlab.com/kode4food/ale/internal/runtime/isa"
)

// Label manages anchoring and referencing of labels
type Label struct {
	encoder  *encoder
	anchored bool
	backRefs []isa.Offset
	offset   isa.Offset
}

const placeholderOffset = isa.Offset(math.MaxUint32)

// NewLabel allocates a Label (
func (e *encoder) NewLabel() *Label {
	return &Label{
		encoder:  e,
		backRefs: []isa.Offset{},
		offset:   placeholderOffset,
	}
}

func (e *encoder) nextOffset() isa.Offset {
	return isa.Offset(len(e.code))
}

// Code turns Label into a Coder, allowing a references to be
// placed at the current encoding position. If the Label has
// not already been anchored, then a pending reference is
// placed until anchoring happens.
func (l *Label) Code() isa.Code {
	e := l.encoder
	if !l.anchored {
		off := e.nextOffset()
		l.backRefs = append(l.backRefs, off)
	}
	return isa.Code(l.offset)
}

// DropAnchor marks the current encoding position as the Label's
// target offset. Any pending references will be finalized
func (l *Label) DropAnchor() {
	if l.anchored {
		panic("label was already anchored")
	}
	e := l.encoder
	off := e.nextOffset()
	l.anchored = true
	l.offset = off
	for _, b := range l.backRefs {
		e.code[b] = isa.Code(off)
	}
	l.backRefs = nil
}
