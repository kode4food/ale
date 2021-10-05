package compound_test

import (
	"testing"

	"github.com/kode4food/ale/types"
	"github.com/kode4food/ale/types/basic"
	"github.com/kode4food/ale/types/compound"
	"github.com/stretchr/testify/assert"
)

func TestApplicable(t *testing.T) {
	as := assert.New(t)

	a1 := compound.Applicable(
		compound.Signature{
			Arguments: []types.Type{basic.Number, basic.Number},
			Result:    basic.Bool,
		},
		compound.Signature{
			Arguments: []types.Type{basic.Number},
			Result:    basic.Bool,
		},
	)

	a2 := compound.Applicable(
		compound.Signature{
			Arguments: []types.Type{basic.Symbol, basic.Bool},
			Result:    basic.Bool,
		},
		compound.Signature{
			Arguments: []types.Type{},
			Result:    basic.Number,
		},
		compound.Signature{
			Arguments: []types.Type{basic.Number},
			Result:    basic.Bool,
		},
		compound.Signature{
			Arguments: []types.Type{basic.Number, basic.Number},
			Result:    basic.Bool,
		},
	)

	as.True(a1.Accepts(a1))
	as.True(a1.Accepts(a2))
	as.False(a2.Accepts(a1))
	as.True(a2.Accepts(a2))

	as.False(a1.Accepts(basic.Number))
	as.False(a1.Accepts(basic.Lambda))
	as.True(basic.Lambda.Accepts(a1))
	as.False(basic.Number.Accepts(a1))

	u1 := compound.Union(a1, a2)
	u2 := compound.Union(compound.List(basic.Symbol), a1)
	as.True(basic.Lambda.Accepts(u1))
	as.False(basic.Lambda.Accepts(u2))
}
