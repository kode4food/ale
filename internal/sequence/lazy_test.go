package sequence_test

import (
	"testing"

	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/internal/sequence"
	"github.com/kode4food/ale/pkg/data"
)

func TestLazySeq(t *testing.T) {
	var inc sequence.LazyResolver
	as := assert.New(t)

	i := 0
	inc = func() (data.Value, data.Sequence, bool) {
		if i >= 10 {
			return data.Null, data.Null, false
		}
		i++
		first := F(float64(i))
		return first, sequence.NewLazy(inc), true
	}

	l := sequence.NewLazy(inc).(data.Prepender).Prepend(F(0))
	as.False(l.IsEmpty())
	as.Number(0, l.Car())
	as.Number(1, l.Cdr().(data.Pair).Car())
	as.Number(2, l.Cdr().(data.Pair).Cdr().(data.Pair).Car())
	as.Contains(":type lazy-sequence", l)
}
