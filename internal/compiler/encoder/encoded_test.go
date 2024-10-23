package encoder_test

import (
	"testing"

	"github.com/kode4food/ale/internal/assert"
	"github.com/kode4food/ale/internal/compiler/encoder"
	"github.com/kode4food/ale/internal/runtime/isa"
)

func makeEncoded(code isa.Instructions) *encoder.Encoded {
	return &encoder.Encoded{
		Code: code,
	}
}

func TestRunnableJump(t *testing.T) {
	as := assert.New(t)

	i1, err := makeEncoded(isa.Instructions{
		isa.NoOp.New(),
		isa.Label.New(0),
		isa.Jump.New(0),
	}).Runnable()
	as.Nil(err)

	as.Equal(isa.Instructions{
		isa.Jump.New(0),
	}, i1.Code)
}

func TestRunnableCondJump(t *testing.T) {
	as := assert.New(t)

	i1, err := makeEncoded(isa.Instructions{
		isa.NoOp.New(),
		isa.Label.New(0),
		isa.NoOp.New(),
		isa.False.New(),
		isa.NoOp.New(),
		isa.CondJump.New(0),
		isa.NoOp.New(),
	}).Runnable()
	as.Nil(err)

	as.Equal(isa.Instructions{
		isa.False.New(),
		isa.CondJump.New(0),
	}, i1.Code)
}

func TestRunnableForwardJump(t *testing.T) {
	as := assert.New(t)

	i1, err := makeEncoded(isa.Instructions{
		isa.NoOp.New(),
		isa.Label.New(0),
		isa.NoOp.New(),
		isa.Jump.New(1),
		isa.NoOp.New(),
		isa.Label.New(1),
		isa.NoOp.New(),
		isa.Jump.New(0),
	}).Runnable()
	as.Nil(err)

	as.Equal(isa.Instructions{
		isa.Jump.New(0),
	}, i1.Code)
}

func TestRunnableDoubleAnchor(t *testing.T) {
	as := assert.New(t)

	_, err := makeEncoded(isa.Instructions{
		isa.Label.New(0),
		isa.Label.New(0),
	}).Runnable()

	as.EqualError(err, encoder.ErrLabelAlreadyAnchored)
}
