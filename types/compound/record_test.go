package compound_test

import (
	"testing"

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

	as.Equal("record", r1.Name())

	as.True(r1.Accepts(r1))
	as.True(r1.Accepts(r2))
	as.False(r2.Accepts(r1))
	as.False(r1.Accepts(r3))
	as.False(r3.Accepts(r1))

	as.False(r1.Accepts(basic.Object))
	as.True(basic.Object.Accepts(r1))
}
