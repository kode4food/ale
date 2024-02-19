package encoder_test

import (
	"testing"

	"github.com/kode4food/ale/compiler/encoder"
	"github.com/kode4food/ale/internal/assert"
	"github.com/kode4food/ale/runtime/isa"
)

func makeEncoded(code isa.Instructions) *encoder.Encoded {
	return &encoder.Encoded{
		Code: code,
	}
}

func TestRunnableJump(t *testing.T) {
	as := assert.New(t)

	i1 := makeEncoded(isa.Instructions{
		isa.NoOp.New(),
		isa.Label.New(0),
		isa.Jump.New(0),
	}).Runnable()

	as.Equal(isa.Instructions{
		isa.Jump.New(0),
	}, i1.Code)
}

func TestRunnableCondJump(t *testing.T) {
	as := assert.New(t)

	i1 := makeEncoded(isa.Instructions{
		isa.NoOp.New(),
		isa.Label.New(0),
		isa.NoOp.New(),
		isa.False.New(),
		isa.NoOp.New(),
		isa.CondJump.New(0),
		isa.NoOp.New(),
	}).Runnable()

	as.Equal(isa.Instructions{
		isa.False.New(),
		isa.CondJump.New(0),
	}, i1.Code)
}

func TestRunnableForwardJump(t *testing.T) {
	as := assert.New(t)

	i1 := makeEncoded(isa.Instructions{
		isa.NoOp.New(),
		isa.Label.New(0),
		isa.NoOp.New(),
		isa.Jump.New(1),
		isa.NoOp.New(),
		isa.Label.New(1),
		isa.NoOp.New(),
		isa.Jump.New(0),
	}).Runnable()

	as.Equal(isa.Instructions{
		isa.Jump.New(0),
	}, i1.Code)
}

func TestRunnableDoubleAnchor(t *testing.T) {
	as := assert.New(t)

	defer as.ExpectPanic(encoder.ErrLabelAlreadyAnchored)

	makeEncoded(isa.Instructions{
		isa.Label.New(0),
		isa.Label.New(0),
	}).Runnable()
}
