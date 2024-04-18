package visitor_test

import (
	"testing"

	"github.com/kode4food/ale/internal/assert"
	"github.com/kode4food/ale/pkg/compiler/encoder"
	"github.com/kode4food/ale/pkg/compiler/generate"
	"github.com/kode4food/ale/pkg/compiler/ir/visitor"
	"github.com/kode4food/ale/pkg/runtime/isa"
)

func TestBranch(t *testing.T) {
	as := assert.New(t)

	e1 := assert.GetTestEncoder()
	e1.Emit(isa.NegInt, 1)
	generate.Branch(e1,
		func(encoder.Encoder) { e1.Emit(isa.True) },
		func(encoder.Encoder) { e1.Emit(isa.PosInt, 1) },
		func(encoder.Encoder) { e1.Emit(isa.Zero) },
	)
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
