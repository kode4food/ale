package visitor_test

import (
	"testing"

	"github.com/kode4food/ale/internal/assert"
	"github.com/kode4food/ale/internal/compiler/encoder"
	"github.com/kode4food/ale/internal/compiler/generate"
	"github.com/kode4food/ale/internal/compiler/ir/visitor"
	"github.com/kode4food/ale/internal/runtime/isa"
)

func TestBranch(t *testing.T) {
	as := assert.New(t)

	e1 := assert.GetTestEncoder()
	e1.Emit(isa.NegInt, 1)
	as.NoError(generate.Branch(e1,
		func(encoder.Encoder) error { e1.Emit(isa.True); return nil },
		func(encoder.Encoder) error { e1.Emit(isa.PosInt, 1); return nil },
		func(encoder.Encoder) error { e1.Emit(isa.Zero); return nil },
	))
	e1.Emit(isa.Pop)
	e1.Emit(isa.Return)

	b := visitor.Branched(e1.Encode().Code).(visitor.Branches)

	as.Instructions(isa.Instructions{
		isa.NegInt.New(1),
		isa.True.New(),
		isa.CondJump.New(0),
	}, b.Prologue().Code())

	as.Instructions(isa.Instructions{
		isa.PosInt.New(1),
	}, b.ThenBranch().Code())

	as.Instructions(isa.Instructions{
		isa.Zero.New(),
	}, b.ElseBranch().Code())

	as.Instructions(isa.Instructions{
		isa.Pop.New(),
		isa.Return.New(),
	}, b.Epilogue().Code())
}
