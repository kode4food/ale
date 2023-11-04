package types_test

import (
	"testing"

	"github.com/kode4food/ale/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestApplicableAccepts(t *testing.T) {
	as := assert.New(t)

	a1 := types.MakeApplicable(
		types.Signature{
			Params: []types.Type{types.BasicNumber, types.BasicNumber},
			Result: types.BasicBoolean,
		},
		types.Signature{
			Params: []types.Type{types.BasicNumber},
			Result: types.BasicBoolean,
		},
	)
	as.Equal(`procedure(number,number->boolean,number->boolean)`, a1.Name())

	a2 := types.MakeApplicable(
		types.Signature{
			Params: []types.Type{types.BasicSymbol, types.BasicBoolean},
			Result: types.BasicBoolean,
		},
		types.Signature{
			Params: []types.Type{},
			Result: types.BasicNumber,
		},
		types.Signature{
			Params: []types.Type{types.BasicNumber},
			Result: types.BasicBoolean,
		},
		types.Signature{
			Params: []types.Type{types.BasicNumber, types.BasicNumber},
			Result: types.BasicBoolean,
		},
	)

	as.Equal(`procedure(symbol,boolean->boolean,->number,number->boolean,number,number->boolean)`, a2.Name())

	as.True(types.Accepts(a1, a1))
	as.True(types.Accepts(a1, a2))
	as.False(types.Accepts(a2, a1))
	as.True(types.Accepts(a2, a2))

	as.False(types.Accepts(a1, types.BasicNumber))
	as.False(types.Accepts(a1, types.BasicProcedure))
	as.True(types.Accepts(types.BasicProcedure, a1))
	as.False(types.Accepts(types.BasicNumber, a1))

	u1 := types.MakeUnion(a1, a2)
	u2 := types.MakeUnion(types.MakeListOf(types.BasicSymbol), a1)
	as.True(types.Accepts(types.BasicProcedure, u1))
	as.False(types.Accepts(types.BasicProcedure, u2))
}

func TestApplicableRest(t *testing.T) {
	as := assert.New(t)

	a1 := types.MakeApplicable(
		types.Signature{
			Params: []types.Type{types.BasicNumber, types.BasicNumber},
			Result: types.BasicBoolean,
		},
	)
	a2 := types.MakeApplicable(
		types.Signature{
			Params:    []types.Type{types.BasicNumber, types.BasicNumber},
			TakesRest: true,
			Result:    types.BasicBoolean,
		},
	)
	a3 := types.MakeApplicable(
		types.Signature{
			Params: []types.Type{types.BasicNumber, types.BasicNumber},
			Result: types.BasicBoolean,
		},
		types.Signature{
			Params:    []types.Type{types.BasicNumber, types.BasicNumber},
			TakesRest: true,
			Result:    types.BasicBoolean,
		},
	)

	as.Equal(`procedure(number.number->boolean)`, a2.Name())

	as.True(types.Accepts(a1, a1))
	as.False(types.Accepts(a1, a2))
	as.False(types.Accepts(a3, a1))
	as.False(types.Accepts(a3, a2))
	as.True(types.Accepts(a1, a3))
	as.True(types.Accepts(a2, a3))
}

func TestApplicableEqual(t *testing.T) {
	as := assert.New(t)

	a1 := types.MakeApplicable(
		types.Signature{
			Params: []types.Type{types.BasicNumber, types.BasicNumber},
			Result: types.BasicBoolean,
		},
	)
	a2 := types.MakeApplicable(
		types.Signature{
			Params:    []types.Type{types.BasicNumber, types.BasicNumber},
			TakesRest: true,
			Result:    types.BasicBoolean,
		},
	)
	a3 := types.MakeApplicable(
		types.Signature{
			Params: []types.Type{types.BasicNumber, types.BasicBoolean},
			Result: types.BasicBoolean,
		},
	)
	a4 := types.MakeApplicable(
		types.Signature{
			Params: []types.Type{types.BasicNumber, types.BasicBoolean},
			Result: types.BasicBoolean,
		},
		types.Signature{
			Params: []types.Type{types.BasicNumber, types.BasicBoolean},
			Result: types.BasicBoolean,
		},
	)
	a5 := types.MakeApplicable(
		types.Signature{
			Params: []types.Type{
				types.BasicNumber, types.BasicBoolean, types.BasicList,
			},
			Result: types.BasicBoolean,
		},
	)

	cpy := *a1
	as.True(a1.Equal(a1))
	as.True(a1.Equal(&cpy))
	as.False(a1.Equal(a2))
	as.False(a2.Equal(types.BasicList))
	as.False(a1.Equal(a3))
	as.False(a3.Equal(a4))
	as.False(a3.Equal(a5))
}
