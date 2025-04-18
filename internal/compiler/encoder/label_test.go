package encoder_test

import (
	"fmt"
	"testing"

	"github.com/kode4food/ale/internal/assert"
	"github.com/kode4food/ale/internal/compiler/ir/analysis"
	"github.com/kode4food/ale/internal/runtime/isa"
)

func TestLabels(t *testing.T) {
	as := assert.New(t)

	e := assert.GetTestEncoder()
	l1 := e.NewLabel()
	l2 := e.NewLabel()
	e.Emit(isa.Jump, l2)
	e.Emit(isa.NoOp)
	e.Emit(isa.Jump, l1)
	e.Emit(isa.Label, l2)
	e.Emit(isa.NoOp)
	e.Emit(isa.Label, l1)

	as.Instructions(isa.Instructions{
		isa.Jump.New(1),
		isa.NoOp.New(),
		isa.Jump.New(0),
		isa.Label.New(1),
		isa.NoOp.New(),
		isa.Label.New(0),
	}, e.Encode().Code)
}

func TestLabelDoubleAnchor(t *testing.T) {
	as := assert.New(t)

	e := assert.GetTestEncoder()
	l1 := e.NewLabel()
	e.Emit(isa.Label, l1)
	e.Emit(isa.Label, l1)
	e.Emit(isa.Jump, l1)

	err := analysis.Verify(e.Encode().Code)
	as.EqualError(err, fmt.Sprintf(analysis.ErrLabelMultipleAnchors, 0))
}
