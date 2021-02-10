package sequence_test

import (
	"testing"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/internal/sequence"
)

func TestLazySeq(t *testing.T) {
	var inc sequence.LazyResolver
	as := assert.New(t)

	i := 0
	inc = func() (data.Value, data.Sequence, bool) {
		if i >= 10 {
			return data.Nil, data.EmptyList, false
		}
		i++
		first := F(float64(i))
		return first, sequence.NewLazy(inc), true
	}

	l := sequence.NewLazy(inc).(data.PrependerSequence).Prepend(F(0))
	as.False(l.IsEmpty())
	as.Number(0, l.First())
	as.Number(1, l.Rest().First())
	as.Number(2, l.Rest().Rest().First())
	as.Contains(":type lazy-sequence", l)
}
