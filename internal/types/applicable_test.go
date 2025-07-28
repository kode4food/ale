package types_test

import (
	"testing"

	"github.com/kode4food/ale"
	"github.com/kode4food/ale/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestApplicableAccepts(t *testing.T) {
	as := assert.New(t)

	a1 := types.MakeApplicable(
		types.Signature{
			Params: []ale.Type{types.BasicNumber, types.BasicNumber},
			Result: types.BasicBoolean,
		},
		types.Signature{
			Params: []ale.Type{types.BasicNumber},
			Result: types.BasicBoolean,
		},
	)
	as.Equal(`procedure(number,number->boolean,number->boolean)`, a1.Name())

	a2 := types.MakeApplicable(
		types.Signature{
			Params: []ale.Type{types.BasicSymbol, types.BasicBoolean},
			Result: types.BasicBoolean,
		},
		types.Signature{
			Params: []ale.Type{},
			Result: types.BasicNumber,
		},
		types.Signature{
			Params: []ale.Type{types.BasicNumber},
			Result: types.BasicBoolean,
		},
		types.Signature{
			Params: []ale.Type{types.BasicNumber, types.BasicNumber},
			Result: types.BasicBoolean,
		},
	)

	as.Equal(
		`procedure(symbol,boolean->boolean,->number,number->boolean,number,number->boolean)`,
		a2.Name(),
	)

	as.True(a1.Accepts(a1))
	as.True(a1.Accepts(a2))
	as.False(a2.Accepts(a1))
	as.True(a2.Accepts(a2))

	as.False(a1.Accepts(types.BasicNumber))
	as.False(a1.Accepts(types.BasicProcedure))
	as.True(types.BasicProcedure.Accepts(a1))
	as.False(types.BasicNumber.Accepts(a1))

	u1 := types.MakeUnion(a1, a2)
	u2 := types.MakeUnion(types.MakeListOf(types.BasicSymbol), a1)
	as.True(types.BasicProcedure.Accepts(u1))
	as.False(types.BasicProcedure.Accepts(u2))
}

func TestApplicableRest(t *testing.T) {
	as := assert.New(t)

	a1 := types.MakeApplicable(
		types.Signature{
			Params: []ale.Type{types.BasicNumber, types.BasicNumber},
			Result: types.BasicBoolean,
		},
	)
	a2 := types.MakeApplicable(
		types.Signature{
			Params:    []ale.Type{types.BasicNumber, types.BasicNumber},
			TakesRest: true,
			Result:    types.BasicBoolean,
		},
	)
	a3 := types.MakeApplicable(
		types.Signature{
			Params: []ale.Type{types.BasicNumber, types.BasicNumber},
			Result: types.BasicBoolean,
		},
		types.Signature{
			Params:    []ale.Type{types.BasicNumber, types.BasicNumber},
			TakesRest: true,
			Result:    types.BasicBoolean,
		},
	)

	as.Equal(`procedure(number.number->boolean)`, a2.Name())

	as.True(a1.Accepts(a1))
	as.False(a1.Accepts(a2))
	as.False(a3.Accepts(a1))
	as.False(a3.Accepts(a2))
	as.True(a1.Accepts(a3))
	as.True(a2.Accepts(a3))
}

func TestApplicableEqual(t *testing.T) {
	as := assert.New(t)

	a1 := types.MakeApplicable(
		types.Signature{
			Params: []ale.Type{types.BasicNumber, types.BasicNumber},
			Result: types.BasicBoolean,
		},
	)
	a2 := types.MakeApplicable(
		types.Signature{
			Params:    []ale.Type{types.BasicNumber, types.BasicNumber},
			TakesRest: true,
			Result:    types.BasicBoolean,
		},
	)
	a3 := types.MakeApplicable(
		types.Signature{
			Params: []ale.Type{types.BasicNumber, types.BasicBoolean},
			Result: types.BasicBoolean,
		},
	)
	a4 := types.MakeApplicable(
		types.Signature{
			Params: []ale.Type{types.BasicNumber, types.BasicBoolean},
			Result: types.BasicBoolean,
		},
		types.Signature{
			Params: []ale.Type{types.BasicNumber, types.BasicBoolean},
			Result: types.BasicBoolean,
		},
	)
	a5 := types.MakeApplicable(
		types.Signature{
			Params: []ale.Type{
				types.BasicNumber, types.BasicBoolean, types.BasicList,
			},
			Result: types.BasicBoolean,
		},
	)

	cpy := *a1.(*types.Applicable)
	as.True(a1.Equal(a1))
	as.True(a1.Equal(&cpy))
	as.False(a1.Equal(a2))
	as.False(a2.Equal(types.BasicList))
	as.False(a1.Equal(a3))
	as.False(a3.Equal(a4))
	as.False(a3.Equal(a5))
}
