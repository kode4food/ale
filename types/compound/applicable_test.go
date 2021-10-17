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
	as.Equal(`lambda(number,number->boolean,number->boolean)`, a1.Name())

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

	as.Equal(`lambda(symbol,boolean->boolean,->number,number->boolean,number,number->boolean)`, a2.Name())

	as.NotNil(types.Check(a1).Accepts(a1))
	as.NotNil(types.Check(a1).Accepts(a2))
	as.Nil(types.Check(a2).Accepts(a1))
	as.NotNil(types.Check(a2).Accepts(a2))

	as.Nil(types.Check(a1).Accepts(basic.Number))
	as.Nil(types.Check(a1).Accepts(basic.Lambda))
	as.NotNil(types.Check(basic.Lambda).Accepts(a1))
	as.Nil(types.Check(basic.Number).Accepts(a1))

	u1 := compound.Union(a1, a2)
	u2 := compound.Union(compound.List(basic.Symbol), a1)
	as.NotNil(types.Check(basic.Lambda).Accepts(u1))
	as.Nil(types.Check(basic.Lambda).Accepts(u2))
}
