package ffi_test

import (
	"testing"

	"github.com/kode4food/ale/core/bootstrap"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/eval"
	"github.com/kode4food/ale/ffi"
	"github.com/kode4food/ale/internal/assert"
)

type (
	EvalWrapped struct {
		*assert.Wrapper
	}

	Env map[data.LocalSymbol]any
)

var testEnv = bootstrap.DevNullEnvironment()

func NewWrapped(t *testing.T) *EvalWrapped {
	return &EvalWrapped{
		Wrapper: assert.New(t),
	}
}

func (e *EvalWrapped) EvalTo(src string, env Env, expect data.Value) {
	e.Helper()
	ns := testEnv.GetAnonymous()
	for n, v := range env {
		v, err := ffi.Wrap(v)
		e.Nil(err)
		ns.Declare(n).Bind(v)
	}
	res := eval.String(ns, data.String(src))
	e.Equal(expect, res)
}
