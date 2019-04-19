package data_test

import (
	"testing"

	"gitlab.com/kode4food/ale/data"
	"gitlab.com/kode4food/ale/internal/assert"
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
