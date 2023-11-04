package data_test

import (
	"testing"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/internal/sequence"
)

func TestLastOfSequence(t *testing.T) {
	as := assert.New(t)

	v, ok := data.Last(data.Null)
	as.Nil(v)
	as.False(ok)

	v, ok = data.Last(L(S("this"), S("is"), S("last")))
	as.String("last", v)
	as.True(ok)

	v, ok = data.Last(V(S("this"), S("is"), S("last")))
	as.String("last", v)
	as.True(ok)

	v, ok = data.Last(sequence.NewLazy(
		func() (data.Value, data.Sequence, bool) {
			return S("hello"), data.Null, true
		},
	))
	as.String("hello", v)
	as.True(ok)

	_, ok = data.Last(sequence.NewLazy(
		func() (data.Value, data.Sequence, bool) {
			return data.Null, data.Null, false
		},
	))
	as.False(ok)
}

func testSequenceCallInterface(as *assert.Wrapper, s data.Procedure) {
	as.NotNil(s.CheckArity(0))
	as.Nil(s.CheckArity(1))
	as.Nil(s.CheckArity(2))
	as.NotNil(s.CheckArity(3))
}
