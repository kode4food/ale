package analysis_test

import (
	"fmt"
	"testing"

	"github.com/kode4food/ale/internal/assert"
	"github.com/kode4food/ale/internal/compiler/encoder"
	"github.com/kode4food/ale/internal/compiler/generate"
	"github.com/kode4food/ale/internal/compiler/ir/analysis"
	"github.com/kode4food/ale/internal/runtime/isa"
)

func TestVerifyGoodStack(t *testing.T) {
	as := assert.New(t)

	e := assert.GetTestEncoder()
	as.Nil(generate.Branch(e,
		func(encoder.Encoder) error { e.Emit(isa.True); return nil },
		func(encoder.Encoder) error { e.Emit(isa.PosInt, 1); return nil },
		func(encoder.Encoder) error { e.Emit(isa.Zero); return nil },
	))
	e.Emit(isa.Return)

	as.Nil(analysis.Verify(e.Encode().Code))
}

func TestVerifyBadBranchStack(t *testing.T) {
	as := assert.New(t)

	e := assert.GetTestEncoder()
	as.Nil(generate.Branch(e,
		func(encoder.Encoder) error { e.Emit(isa.True); return nil },
		func(encoder.Encoder) error {
			e.Emit(isa.PosInt, 1)
			e.Emit(isa.PosInt, 2)
			return nil
		},
		func(encoder.Encoder) error { e.Emit(isa.Zero); return nil },
	))
	e.Emit(isa.Return)

	err := analysis.Verify(e.Encode().Code)
	as.EqualError(err, analysis.ErrBadBranchTermination)
}

func TestVerifyBadEndStack(t *testing.T) {
	as := assert.New(t)

	e := assert.GetTestEncoder()
	e.Emit(isa.True)

	err := analysis.Verify(e.Encode().Code)
	as.EqualError(err, fmt.Sprintf(analysis.ErrBadStackTermination, 1))
}

func TestVerifyGoodJump(t *testing.T) {
	as := assert.New(t)

	e := assert.GetTestEncoder()
	lbl := e.NewLabel()
	e.Emit(isa.Label, lbl)
	e.Emit(isa.True)
	e.Emit(isa.Pop)
	e.Emit(isa.Jump, lbl)

	as.Nil(analysis.Verify(e.Encode().Code))
}

func TestVerifyBadJump(t *testing.T) {
	as := assert.New(t)

	e := assert.GetTestEncoder()
	lbl := e.NewLabel()
	e.Emit(isa.True)
	e.Emit(isa.Pop)
	e.Emit(isa.Jump, lbl)

	err := analysis.Verify(e.Encode().Code)
	as.EqualError(err, fmt.Sprintf(analysis.ErrLabelNotAnchored, 0))
}
