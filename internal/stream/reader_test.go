package stream_test

import (
	"bytes"
	"testing"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/assert"
	"github.com/kode4food/ale/internal/stream"
)

func TestReader(t *testing.T) {
	as := assert.New(t)

	b1 := []byte("12ö34")
	b2 := []byte("12\n34")

	r1 := stream.NewReader(bytes.NewReader(b1), stream.RuneInput)
	as.False(r1.IsEmpty())
	as.String("1", r1.Car())
	as.String("2", r1.Cdr().(data.Pair).Car())
	as.String("ö", r1.Cdr().(data.Pair).Cdr().(data.Pair).Car())
	as.String("3", r1.Cdr().(data.Pair).Cdr().(data.Pair).
		Cdr().(data.Pair).Car())
	as.String("4", r1.Cdr().(data.Pair).Cdr().(data.Pair).
		Cdr().(data.Pair).Cdr().(data.Pair).Car())
	s1 := r1.Cdr().(data.Pair).Cdr().(data.Pair).Cdr().(data.Pair).
		Cdr().(data.Pair).Cdr().(data.Sequence)
	as.True(s1.IsEmpty())

	r2 := stream.NewReader(bytes.NewReader(b1), stream.LineInput)
	as.False(r2.IsEmpty())
	as.String("12ö34", r2.Car())
	as.True(r2.Cdr().(data.Sequence).IsEmpty())

	r3 := stream.NewReader(bytes.NewReader(b2), stream.LineInput)
	as.False(r3.IsEmpty())
	as.String("12", r3.Car())
	as.String("34", r3.Cdr().(data.Pair).Car())
	as.True(r3.Cdr().(data.Pair).Cdr().(data.Sequence).IsEmpty())
}
