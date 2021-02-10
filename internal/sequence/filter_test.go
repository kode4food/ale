package sequence_test

import (
	"testing"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/internal/sequence"
)

func TestFilter(t *testing.T) {
	as := assert.New(t)

	filterTest := data.Applicative(func(args ...data.Value) data.Value {
		return B(string(args[0].(data.String)) != "filtered out")
	}, 1)

	l := L(S("first"), S("filtered out"), S("last"))
	w := sequence.Filter(l, filterTest)

	v1 := w.First()
	as.String("first", v1)

	v2 := w.Rest().First()
	as.String("last", v2)

	r1 := w.Rest().Rest()
	as.True(r1.IsEmpty())

	p := w.(data.PrependerSequence).Prepend(S("filtered out"))
	v4 := p.First()
	r2 := p.Rest()
	as.String("filtered out", v4)
	as.Equal(w.First(), r2.First())
}

func TestFiltered(t *testing.T) {
	as := assert.New(t)

	l := L(S("first"), S("middle"), S("last"))
	fn1 := data.Applicative(func(args ...data.Value) data.Value {
		return B(string(args[0].(data.String)) != "middle")
	}, 1)
	w1 := sequence.Filter(l, fn1)
	v1 := w1.First()
	as.String("first", v1)

	v2 := w1.Rest().First()
	as.String("last", v2)

	r1 := w1.Rest().Rest()
	as.True(r1.IsEmpty())
}
