package encoder

import "gitlab.com/kode4food/ale/internal/runtime/isa"

type (
	// Label manages anchoring and referencing of labels
	Label struct {
		encoder  *encoder
		number   isa.Index
		anchored bool
	}
)

// NewLabel allocates a Label (
func (e *encoder) NewLabel() *Label {
	res := &Label{
		encoder: e,
		number:  isa.Index(e.nextLabel),
	}
	e.nextLabel++
	return res
}

// Word turns Label into a Coder, allowing a references to be
// placed at the current encoding position.
func (l *Label) Word() isa.Word {
	return isa.Word(l.number)
}

// DropAnchor marks the current encoding position as the Label target
func (l *Label) DropAnchor() {
	if l.anchored {
		panic("label has already been anchored")
	}
	e := l.encoder
	e.Append(isa.Label, l.number)
	l.anchored = true
}
