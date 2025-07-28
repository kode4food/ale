package types_test

import (
	"testing"

	"github.com/kode4food/ale/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestRecordAccepts(t *testing.T) {
	as := assert.New(t)

	r1 := types.MakeRecord(
		types.Field{
			Name:  "first",
			Value: types.BasicString,
		},
		types.Field{
			Name:  "last",
			Value: types.BasicString,
		},
	)

	r2 := types.MakeRecord(
		types.Field{
			Name:  "age",
			Value: types.BasicNumber,
		},
		types.Field{
			Name:  "first",
			Value: types.BasicString,
		},
		types.Field{
			Name:  "last",
			Value: types.BasicString,
		},
	)

	r3 := types.MakeRecord(
		types.Field{
			Name:  "first",
			Value: types.BasicString,
		},
		types.Field{
			Name:  "last",
			Value: types.BasicKeyword,
		},
	)

	as.Equal(`record("first"->string,"last"->keyword)`, r3.Name())

	as.True(r1.Accepts(r1))
	as.True(r1.Accepts(r2))
	as.False(r2.Accepts(r1))
	as.False(r1.Accepts(r3))
	as.False(r3.Accepts(r1))

	as.False(r1.Accepts(types.BasicObject))
	as.True(types.BasicObject.Accepts(r1))
}

func TestRecordEqual(t *testing.T) {
	as := assert.New(t)

	r1 := types.MakeRecord(
		types.Field{
			Name:  "first",
			Value: types.BasicString,
		},
		types.Field{
			Name:  "last",
			Value: types.BasicString,
		},
	)

	r2 := types.MakeRecord(
		types.Field{
			Name:  "age",
			Value: types.BasicNumber,
		},
		types.Field{
			Name:  "first",
			Value: types.BasicString,
		},
		types.Field{
			Name:  "last",
			Value: types.BasicString,
		},
	)

	r3 := types.MakeRecord(
		types.Field{
			Name:  "first",
			Value: types.BasicString,
		},
		types.Field{
			Name:  "last",
			Value: types.BasicKeyword,
		},
	)

	cpy := *r1.(*types.Record)
	as.True(r1.Equal(r1))
	as.True(r1.Equal(&cpy))
	as.False(r1.Equal(r2))
	as.False(r2.Equal(r1))
	as.False(r1.Equal(r3))
	as.False(r3.Equal(r1))

	as.False(r1.Equal(types.BasicObject))
	as.False(types.BasicObject.Equal(r1))
}

func TestRecordNameEscape(t *testing.T) {
	as := assert.New(t)

	r1 := types.MakeRecord(
		types.Field{
			Name:  `I am "quoted"`,
			Value: types.BasicString,
		},
	)

	as.Equal(`record("I am \"quoted\""->string)`, r1.Name())
}
