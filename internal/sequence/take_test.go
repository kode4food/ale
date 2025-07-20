package sequence_test

import (
	"testing"

	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/internal/sequence"
	"github.com/kode4food/ale/data"
)

func TestTake(t *testing.T) {
	as := assert.New(t)
	l := L(I(1), I(2), I(3), I(4))

	v, r, ok := sequence.Take(l, 3)
	as.Equal(V(I(1), I(2), I(3)), v)
	as.Equal(L(I(4)), r)
	as.True(ok)

	v, r, ok = sequence.Take(l, 4)
	as.Equal(V(I(1), I(2), I(3), I(4)), v)
	as.Equal(data.Null, r)
	as.True(ok)

	v, r, ok = sequence.Take(l, 5)
	as.Equal(data.EmptyVector, v)
	as.Equal(data.Null, r)
	as.False(ok)
}
