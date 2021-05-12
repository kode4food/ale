package encoder_test

import (
	"testing"

	"github.com/kode4food/ale/compiler/encoder"
	"github.com/kode4food/ale/core/bootstrap"
	"github.com/kode4food/ale/env"
	"github.com/kode4food/ale/internal/assert"
	"github.com/kode4food/ale/runtime/isa"
)

type EvalWrapped struct {
	*assert.Wrapper
}

var (
	testEnv = env.NewEnvironment()
	ready   bool
)

func getTestEncoder() encoder.Encoder {
	if !ready {
		bootstrap.Into(testEnv)
		ready = true
	}
	ns := testEnv.GetAnonymous()
	return encoder.NewEncoder(ns)
}

func NewWrapped(t *testing.T) *EvalWrapped {
	return &EvalWrapped{
		Wrapper: assert.New(t),
	}
}

func (e *EvalWrapped) Instructions(expected, actual isa.Instructions) {
	e.Equal(len(expected), len(actual))
	for i, l := range expected {
		e.Equal(l, actual[i])
	}
}
