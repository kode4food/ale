package types_test

import (
	"testing"

	"github.com/kode4food/ale/types"
	"github.com/stretchr/testify/assert"
)

func TestRecord(t *testing.T) {
	as := assert.New(t)

	r1 := types.Record(
		types.Field{
			Name:  "first",
			Value: types.String,
		},
		types.Field{
			Name:  "last",
			Value: types.String,
		},
	)

	r2 := types.Record(
		types.Field{
			Name:  "age",
			Value: types.Number,
		},
		types.Field{
			Name:  "first",
			Value: types.String,
		},
		types.Field{
			Name:  "last",
			Value: types.String,
		},
	)

	r3 := types.Record(
		types.Field{
			Name:  "first",
			Value: types.String,
		},
		types.Field{
			Name:  "last",
			Value: types.Keyword,
		},
	)

	as.Equal(`record("first"->string,"last"->keyword)`, r3.Name())

	as.True(types.Accepts(r1, r1))
	as.True(types.Accepts(r1, r2))
	as.False(types.Accepts(r2, r1))
	as.False(types.Accepts(r1, r3))
	as.False(types.Accepts(r3, r1))

	as.False(types.Accepts(r1, types.AnyObject))
	as.True(types.Accepts(types.AnyObject, r1))
}

func TestRecordNameEscape(t *testing.T) {
	as := assert.New(t)

	r1 := types.Record(
		types.Field{
			Name:  `I am "quoted"`,
			Value: types.String,
		},
	)

	as.Equal(`record("I am \"quoted\""->string)`, r1.Name())
}
