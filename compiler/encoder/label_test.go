package encoder_test

import (
	"testing"

	"github.com/kode4food/ale/compiler/ir/analysis"
	"github.com/kode4food/ale/internal/assert"
	"github.com/kode4food/ale/runtime/isa"
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
		isa.New(isa.Jump, 1),
		isa.New(isa.NoOp),
		isa.New(isa.Jump, 0),
		isa.New(isa.Label, 1),
		isa.New(isa.NoOp),
		isa.New(isa.Label, 0),
	}, e.Code())
}

func TestLabelDoubleAnchor(t *testing.T) {
	as := assert.New(t)

	e := assert.GetTestEncoder()
	l1 := e.NewLabel()
	e.Emit(isa.Label, l1)
	e.Emit(isa.Label, l1)
	e.Emit(isa.Jump, l1)

	defer as.ExpectPanic("label anchored multiple times: 0")
	analysis.Verify(e.Code())
}
