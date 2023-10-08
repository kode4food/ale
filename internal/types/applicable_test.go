package types_test

import (
	"testing"

	"github.com/kode4food/ale/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestApplicable(t *testing.T) {
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
	as.Equal(`lambda(number,number->boolean,number->boolean)`, a1.Name())

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

	as.Equal(`lambda(symbol,boolean->boolean,->number,number->boolean,number,number->boolean)`, a2.Name())

	as.True(types.Accepts(a1, a1))
	as.True(types.Accepts(a1, a2))
	as.False(types.Accepts(a2, a1))
	as.True(types.Accepts(a2, a2))

	as.False(types.Accepts(a1, types.BasicNumber))
	as.False(types.Accepts(a1, types.BasicLambda))
	as.True(types.Accepts(types.BasicLambda, a1))
	as.False(types.Accepts(types.BasicNumber, a1))

	u1 := types.MakeUnion(a1, a2)
	u2 := types.MakeUnion(types.MakeListOf(types.BasicSymbol), a1)
	as.True(types.Accepts(types.BasicLambda, u1))
	as.False(types.Accepts(types.BasicLambda, u2))
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

	as.False(types.Accepts(a1, a2))
	as.False(types.Accepts(a3, a1))
	as.False(types.Accepts(a3, a2))
	as.True(types.Accepts(a1, a3))
	as.True(types.Accepts(a2, a3))
}
