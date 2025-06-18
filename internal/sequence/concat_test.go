package sequence_test

import (
	"testing"

	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/internal/sequence"
	"github.com/kode4food/ale/pkg/data"
)

func TestConcat(t *testing.T) {
	as := assert.New(t)

	c0 := sequence.Concat()
	as.Equal(data.Null, c0)

	l1 := L(I(1), I(2), I(3))
	l2 := L(I(4), I(5), I(6))
	l3 := L(I(7), I(8), I(9))

	c1 := sequence.Concat(l1)
	as.True(c1 == l1)
	as.Equal(
		V(I(1), I(2), I(3)),
		sequence.ToVector(c1),
	)

	c2 := sequence.Concat(l1, l2, l3)
	as.Equal(
		V(I(1), I(2), I(3), I(4), I(5), I(6), I(7), I(8), I(9)),
		sequence.ToVector(c2),
	)

	c3 := sequence.Concat(l1, l2)
	_, r, _ := c3.Split()
	_, r, _ = r.Split()
	_, r, _ = r.Split()
	_, r, _ = r.Split()
	as.Equal(l2.Cdr(), r)

	as.Equal(
		V(I(1), I(2), I(3), I(4), I(5), I(6)),
		sequence.ToVector(c3),
	)
}
