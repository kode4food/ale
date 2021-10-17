package compound_test

import (
	"testing"

	"github.com/kode4food/ale/types"
	"github.com/kode4food/ale/types/basic"
	"github.com/kode4food/ale/types/compound"
	"github.com/stretchr/testify/assert"
)

func TestRecord(t *testing.T) {
	as := assert.New(t)

	r1 := compound.Record(
		compound.Field{
			Name:  "first",
			Value: basic.String,
		},
		compound.Field{
			Name:  "last",
			Value: basic.String,
		},
	)

	r2 := compound.Record(
		compound.Field{
			Name:  "age",
			Value: basic.Number,
		},
		compound.Field{
			Name:  "first",
			Value: basic.String,
		},
		compound.Field{
			Name:  "last",
			Value: basic.String,
		},
	)

	r3 := compound.Record(
		compound.Field{
			Name:  "first",
			Value: basic.String,
		},
		compound.Field{
			Name:  "last",
			Value: basic.Keyword,
		},
	)

	as.Equal(`record("first"->string,"last"->keyword)`, r3.Name())

	as.NotNil(types.Check(r1).Accepts(r1))
	as.NotNil(types.Check(r1).Accepts(r2))
	as.Nil(types.Check(r2).Accepts(r1))
	as.Nil(types.Check(r1).Accepts(r3))
	as.Nil(types.Check(r3).Accepts(r1))

	as.Nil(types.Check(r1).Accepts(basic.Object))
	as.NotNil(types.Check(basic.Object).Accepts(r1))
}

func TestRecordNameEscape(t *testing.T) {
	as := assert.New(t)

	r1 := compound.Record(
		compound.Field{
			Name:  `I am "quoted"`,
			Value: basic.String,
		},
	)

	as.Equal(`record("I am \"quoted\""->string)`, r1.Name())
}
