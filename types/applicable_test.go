package types_test

import (
	"testing"

	"github.com/kode4food/ale/types"
	"github.com/stretchr/testify/assert"
)

func TestApplicable(t *testing.T) {
	as := assert.New(t)

	a1 := types.Applicable(
		types.Signature{
			Arguments: []types.Type{types.Number, types.Number},
			Result:    types.Bool,
		},
		types.Signature{
			Arguments: []types.Type{types.Number},
			Result:    types.Bool,
		},
	)
	as.Equal(`lambda(number,number->boolean,number->boolean)`, a1.Name())

	a2 := types.Applicable(
		types.Signature{
			Arguments: []types.Type{types.Symbol, types.Bool},
			Result:    types.Bool,
		},
		types.Signature{
			Arguments: []types.Type{},
			Result:    types.Number,
		},
		types.Signature{
			Arguments: []types.Type{types.Number},
			Result:    types.Bool,
		},
		types.Signature{
			Arguments: []types.Type{types.Number, types.Number},
			Result:    types.Bool,
		},
	)

	as.Equal(`lambda(symbol,boolean->boolean,->number,number->boolean,number,number->boolean)`, a2.Name())

	as.True(types.Accepts(a1, a1))
	as.True(types.Accepts(a1, a2))
	as.False(types.Accepts(a2, a1))
	as.True(types.Accepts(a2, a2))

	as.False(types.Accepts(a1, types.Number))
	as.False(types.Accepts(a1, types.Lambda))
	as.True(types.Accepts(types.Lambda, a1))
	as.False(types.Accepts(types.Number, a1))

	u1 := types.Union(a1, a2)
	u2 := types.Union(types.ListOf(types.Symbol), a1)
	as.True(types.Accepts(types.Lambda, u1))
	as.False(types.Accepts(types.Lambda, u2))
}

func TestApplicableRest(t *testing.T) {
	as := assert.New(t)

	a1 := types.Applicable(
		types.Signature{
			Arguments: []types.Type{types.Number, types.Number},
			Result:    types.Bool,
		},
	)
	a2 := types.Applicable(
		types.Signature{
			Arguments: []types.Type{types.Number, types.Number},
			TakesRest: true,
			Result:    types.Bool,
		},
	)
	a3 := types.Applicable(
		types.Signature{
			Arguments: []types.Type{types.Number, types.Number},
			Result:    types.Bool,
		},
		types.Signature{
			Arguments: []types.Type{types.Number, types.Number},
			TakesRest: true,
			Result:    types.Bool,
		},
	)

	as.False(types.Accepts(a1, a2))
	as.False(types.Accepts(a3, a1))
	as.False(types.Accepts(a3, a2))
	as.True(types.Accepts(a1, a3))
	as.True(types.Accepts(a2, a3))
}
