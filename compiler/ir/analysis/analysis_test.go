package analysis_test

import (
	"testing"

	"github.com/kode4food/ale/compiler/generate"
	"github.com/kode4food/ale/compiler/ir/analysis"
	"github.com/kode4food/ale/internal/assert"
	"github.com/kode4food/ale/runtime/isa"
)

func TestVerifyGoodStack(t *testing.T) {
	as := assert.New(t)

	e := assert.GetTestEncoder()
	generate.Branch(e,
		func() { e.Emit(isa.True) },
		func() { e.Emit(isa.One) },
		func() { e.Emit(isa.Zero) },
	)
	e.Emit(isa.Return)

	defer as.ExpectNoPanic()
	analysis.Verify(e.Code())
}

func TestVerifyBadBranchStack(t *testing.T) {
	as := assert.New(t)

	e := assert.GetTestEncoder()
	generate.Branch(e,
		func() { e.Emit(isa.True) },
		func() {
			e.Emit(isa.One)
			e.Emit(isa.Two)
		},
		func() { e.Emit(isa.Zero) },
	)
	e.Emit(isa.Return)

	defer as.ExpectProgrammerError("branches should end in the same state")
	analysis.Verify(e.Code())
}

func TestVerifyBadEndStack(t *testing.T) {
	as := assert.New(t)

	e := assert.GetTestEncoder()
	e.Emit(isa.True)

	defer as.ExpectProgrammerError("invalid stack end-state: 1")
	analysis.Verify(e.Code())
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
	analysis.Verify(e.Code())
}

func TestVerifyBadJump(t *testing.T) {
	as := assert.New(t)

	e := assert.GetTestEncoder()
	lbl := e.NewLabel()
	e.Emit(isa.True)
	e.Emit(isa.Pop)
	e.Emit(isa.Jump, lbl)

	defer as.ExpectPanic("label not anchored: 0")
	analysis.Verify(e.Code())
}
