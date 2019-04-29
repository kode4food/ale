package stdlib_test

import (
	"testing"

	"gitlab.com/kode4food/ale/data"
	"gitlab.com/kode4food/ale/internal/assert"
	. "gitlab.com/kode4food/ale/internal/assert/helpers"
	"gitlab.com/kode4food/ale/stdlib"
)

func TestLazySeq(t *testing.T) {
	var inc stdlib.LazyResolver
	as := assert.New(t)

	i := 0
	inc = func() (data.Value, data.Sequence, bool) {
		if i >= 10 {
			return data.Nil, data.EmptyList, false
		}
		i++
		first := F(float64(i))
		return first, stdlib.NewLazySequence(inc), true
	}

	l := stdlib.NewLazySequence(inc).Prepend(F(0))
	as.False(l.IsEmpty())
	as.Number(0, l.First())
	as.Number(1, l.Rest().First())
	as.Number(2, l.Rest().Rest().First())
	as.Contains(":type lazy-sequence", l)
}
