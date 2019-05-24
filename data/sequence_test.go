package data_test

import (
	"testing"

	"gitlab.com/kode4food/ale/data"
	"gitlab.com/kode4food/ale/internal/assert"
	. "gitlab.com/kode4food/ale/internal/assert/helpers"
)

type ncSeq struct{}

func (n *ncSeq) First() data.Value {
	return nil
}

func (n *ncSeq) Rest() data.Sequence {
	return nil
}

func (n *ncSeq) Split() (data.Value, data.Sequence, bool) {
	return nil, nil, false
}

func (n *ncSeq) Prepend(v data.Value) data.Sequence {
	return nil
}

func (n *ncSeq) IsEmpty() bool {
	return true
}

func (n *ncSeq) String() string {
	return "()"
}

func TestLastOfSequence(t *testing.T) {
	as := assert.New(t)

	v, ok := data.Last(data.EmptyList)
	as.Nil(v)
	as.False(ok)

	v, ok = data.Last(L(S("this"), S("is"), S("last")))
	as.String("last", v)
	as.True(ok)

	v, ok = data.Last(V(S("this"), S("is"), S("last")))
	as.String("last", v)
	as.True(ok)
}
