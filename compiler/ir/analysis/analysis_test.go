package analysis_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/kode4food/ale/compiler/encoder"
	"github.com/kode4food/ale/compiler/generate"
	"github.com/kode4food/ale/compiler/ir/analysis"
	"github.com/kode4food/ale/internal/assert"
	"github.com/kode4food/ale/runtime/isa"
)

func TestVerifyGoodStack(t *testing.T) {
	as := assert.New(t)

	e := assert.GetTestEncoder()
	generate.Branch(e,
		func(encoder.Encoder) { e.Emit(isa.True) },
		func(encoder.Encoder) { e.Emit(isa.PosInt, 1) },
		func(encoder.Encoder) { e.Emit(isa.Zero) },
	)
	e.Emit(isa.Return)

	defer as.ExpectNoPanic()
	analysis.Verify(e.Encode().Code)
}

func TestVerifyBadBranchStack(t *testing.T) {
	as := assert.New(t)

	e := assert.GetTestEncoder()
	generate.Branch(e,
		func(encoder.Encoder) { e.Emit(isa.True) },
		func(encoder.Encoder) {
			e.Emit(isa.PosInt, 1)
			e.Emit(isa.PosInt, 2)
		},
		func(encoder.Encoder) { e.Emit(isa.Zero) },
	)
	e.Emit(isa.Return)

	defer as.ExpectPanic(errors.New(analysis.ErrBadBranchTermination))
	analysis.Verify(e.Encode().Code)
}

func TestVerifyBadEndStack(t *testing.T) {
	as := assert.New(t)

	e := assert.GetTestEncoder()
	e.Emit(isa.True)

	defer as.ExpectPanic(fmt.Errorf(analysis.ErrBadStackTermination, 1))
	analysis.Verify(e.Encode().Code)
}

func TestVerifyGoodJump(t *testing.T) {
	as := assert.New(t)

	e := assert.GetTestEncoder()
	lbl := e.NewLabel()
	e.Emit(isa.Label, lbl)
	e.Emit(isa.True)
	e.Emit(isa.Pop)
	e.Emit(isa.Jump, lbl)

	defer as.ExpectNoPanic()
	analysis.Verify(e.Encode().Code)
}

func TestVerifyBadJump(t *testing.T) {
	as := assert.New(t)

	e := assert.GetTestEncoder()
	lbl := e.NewLabel()
	e.Emit(isa.True)
	e.Emit(isa.Pop)
	e.Emit(isa.Jump, lbl)

	defer as.ExpectPanic(fmt.Errorf(analysis.ErrLabelNotAnchored, 0))
	analysis.Verify(e.Encode().Code)
}
