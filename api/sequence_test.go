package api_test

import (
	"testing"

	"gitlab.com/kode4food/ale/api"
	"gitlab.com/kode4food/ale/internal/assert"
)

type ncSeq struct{}

func (n *ncSeq) First() api.Value                       { return nil }
func (n *ncSeq) Rest() api.Sequence                     { return nil }
func (n *ncSeq) Split() (api.Value, api.Sequence, bool) { return nil, nil, false }
func (n *ncSeq) Prepend(v api.Value) api.Sequence       { return nil }
func (n *ncSeq) IsSequence() bool                       { return false }
func (n *ncSeq) String() string                         { return "()" }

func TestNonCountableSequence(t *testing.T) {
	as := assert.New(t)
	nc := &ncSeq{}

	e := cvtErr("*api_test.ncSeq", "api.CountedSequence", "Count")
	defer as.ExpectPanic(e)
	api.Count(nc)
}
