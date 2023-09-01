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

	v1 := w.Car()
	as.String("first", v1)

	v2 := w.Cdr().(data.Pair).Car()
	as.String("last", v2)

	r1 := w.Cdr().(data.Pair).Cdr()
	as.True(r1.(data.Sequence).IsEmpty())

	p := w.(data.Prepender).Prepend(S("filtered out"))
	v4 := p.Car()
	r2 := p.Cdr()
	as.String("filtered out", v4)
	as.Equal(w.Car(), r2.(data.Pair).Car())
}

func TestFiltered(t *testing.T) {
	as := assert.New(t)

	l := L(S("first"), S("middle"), S("last"))
	fn1 := data.Applicative(func(args ...data.Value) data.Value {
		return B(string(args[0].(data.String)) != "middle")
	}, 1)
	w1 := sequence.Filter(l, fn1)
	v1 := w1.Car()
	as.String("first", v1)

	v2 := w1.Cdr().(data.Pair).Car()
	as.String("last", v2)

	r1 := w1.Cdr().(data.Pair).Cdr()
	as.True(r1.(data.Sequence).IsEmpty())
}
