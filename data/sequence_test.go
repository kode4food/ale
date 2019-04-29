package data_test

import (
	"testing"

	"gitlab.com/kode4food/ale/stdlib"

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

func TestNonCountableSequence(t *testing.T) {
	as := assert.New(t)
	nc := &ncSeq{}

	e := cvtErr("*data_test.ncSeq", "data.CountedSequence", "Count")
	defer as.ExpectPanic(e)
	data.Count(nc)
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

func TestLazyLastOfSequence(t *testing.T) {
	as := assert.New(t)

	v1 := V(I(1), I(2), I(3))
	l1 := stdlib.Map(v1, func(args ...data.Value) data.Value {
		return args[0].(data.Integer) * 2
	})

	v, ok := data.Last(l1)
	as.Number(6, v)
	as.True(ok)
}
