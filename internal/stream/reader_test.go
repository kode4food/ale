package stream_test

import (
	"bytes"
	"testing"

	"github.com/kode4food/ale/internal/assert"
	"github.com/kode4food/ale/internal/stream"
)

func TestReader(t *testing.T) {
	as := assert.New(t)

	b1 := []byte("12ö34")
	b2 := []byte("12\n34")

	r1 := stream.NewReader(bytes.NewReader(b1), stream.RuneInput)
	as.False(r1.IsEmpty())
	as.String("1", r1.First())
	as.String("2", r1.Rest().First())
	as.String("ö", r1.Rest().Rest().First())
	as.String("3", r1.Rest().Rest().Rest().First())
	as.String("4", r1.Rest().Rest().Rest().Rest().First())
	s1 := r1.Rest().Rest().Rest().Rest().Rest()
	as.True(s1.IsEmpty())

	r2 := stream.NewReader(bytes.NewReader(b1), stream.LineInput)
	as.False(r2.IsEmpty())
	as.String("12ö34", r2.First())
	as.True(r2.Rest().IsEmpty())

	r3 := stream.NewReader(bytes.NewReader(b2), stream.LineInput)
	as.False(r3.IsEmpty())
	as.String("12", r3.First())
	as.String("34", r3.Rest().First())
	as.True(r3.Rest().Rest().IsEmpty())
}
