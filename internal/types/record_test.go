package types_test

import (
	"testing"

	"github.com/kode4food/ale/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestRecord(t *testing.T) {
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

	as.True(types.Accepts(r1, r1))
	as.True(types.Accepts(r1, r2))
	as.False(types.Accepts(r2, r1))
	as.False(types.Accepts(r1, r3))
	as.False(types.Accepts(r3, r1))

	as.False(types.Accepts(r1, types.BasicObject))
	as.True(types.Accepts(types.BasicObject, r1))
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
