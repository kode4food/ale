package visitor_test

import (
	"fmt"
	"testing"

	"github.com/kode4food/ale/compiler/ir/visitor"

	"github.com/kode4food/ale/compiler/generate"

	"github.com/kode4food/ale/internal/assert"
	"github.com/kode4food/ale/runtime/isa"
)

func TestBranch(t *testing.T) {
	as := assert.New(t)

	e1 := assert.GetTestEncoder()
	e1.Emit(isa.NegOne)
	generate.Branch(e1,
		func() { e1.Emit(isa.True) },
		func() { e1.Emit(isa.One) },
		func() { e1.Emit(isa.Zero) },
	)
	e1.Emit(isa.Pop)
	e1.Emit(isa.Return)

	b := visitor.Branch(e1.Code()).(visitor.Branches)

	as.Instructions(isa.Instructions{
		isa.New(isa.NegOne),
		isa.New(isa.True),
		isa.New(isa.CondJump, 0),
	}, b.Prologue().Code())

	as.Instructions(isa.Instructions{
		isa.New(isa.One),
	}, b.ThenBranch().Code())

	as.Instructions(isa.Instructions{
		isa.New(isa.Zero),
	}, b.ElseBranch().Code())

	as.Instructions(isa.Instructions{
		isa.New(isa.Pop),
		isa.New(isa.Return),
	}, b.Epilogue().Code())

	as.Equal(`NegOne()
True()
CondJump(0)
  One()
else:
  Zero()
Pop()
Return()
`, b.(fmt.Stringer).String())

	as.Equal(`Pop()
Return()
`, b.Epilogue().(fmt.Stringer).String())
}
