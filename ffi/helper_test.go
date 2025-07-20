package ffi_test

import (
	"testing"

	"github.com/kode4food/ale"
	"github.com/kode4food/ale/core/bootstrap"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/env"
	"github.com/kode4food/ale/eval"
	"github.com/kode4food/ale/ffi"
	"github.com/kode4food/ale/internal/assert"
)

type (
	EvalWrapped struct {
		*assert.Wrapper
	}

	Env map[data.Local]any
)

var testEnv = bootstrap.DevNullEnvironment()

func NewWrapped(t *testing.T) *EvalWrapped {
	return &EvalWrapped{
		Wrapper: assert.New(t),
	}
}

func (e *EvalWrapped) EvalTo(src string, en Env, expect ale.Value) {
	e.Helper()
	ns := testEnv.GetAnonymous()
	for n, v := range en {
		v, err := ffi.Wrap(v)
		if e.NoError(err) {
			e.NoError(env.BindPublic(ns, n, v))
		}
	}
	res, err := eval.String(ns, data.String(src))
	if e.NoError(err) {
		e.Equal(expect, res)
	}
}
