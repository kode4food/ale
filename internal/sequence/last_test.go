package sequence_test

import (
	"testing"

	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/internal/sequence"
	"github.com/kode4food/ale/pkg/data"
)

func TestLast(t *testing.T) {
	as := assert.New(t)

	v, ok := sequence.Last(data.Null)
	as.Nil(v)
	as.False(ok)

	v, ok = sequence.Last(L(S("this"), S("is"), S("last")))
	as.String("last", v)
	as.True(ok)

	v, ok = sequence.Last(V(S("this"), S("is"), S("last")))
	as.String("last", v)
	as.True(ok)

	v, ok = sequence.Last(sequence.NewLazy(
		func() (data.Value, data.Sequence, bool) {
			return S("hello"), data.Null, true
		},
	))
	as.String("hello", v)
	as.True(ok)

	_, ok = sequence.Last(sequence.NewLazy(
		func() (data.Value, data.Sequence, bool) {
			return data.Null, data.Null, false
		},
	))
	as.False(ok)
}
