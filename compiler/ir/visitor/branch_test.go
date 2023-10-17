package visitor_test

import (
	"fmt"
	"testing"

	"github.com/kode4food/ale/compiler/encoder"
	"github.com/kode4food/ale/compiler/generate"
	"github.com/kode4food/ale/compiler/ir/visitor"
	"github.com/kode4food/ale/internal/assert"
	"github.com/kode4food/ale/runtime/isa"
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

	b := visitor.Branch(e1.Code()).(visitor.Branches)

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

	as.Equal(
		"neg-int 1\ntrue\ncond-jump 0\n  pos-int 1\nelse:\n  zero\npop\nreturn\n",
		b.(fmt.Stringer).String(),
	)

	as.Equal("pop\nreturn\n", b.Epilogue().(fmt.Stringer).String())
}
