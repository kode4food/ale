package ffi_test

import (
	"testing"

	"github.com/kode4food/ale/internal/assert"
	"github.com/kode4food/ale/pkg/core/bootstrap"
	"github.com/kode4food/ale/pkg/data"
	"github.com/kode4food/ale/pkg/env"
	"github.com/kode4food/ale/pkg/eval"
	"github.com/kode4food/ale/pkg/ffi"
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

func (e *EvalWrapped) EvalTo(src string, en Env, expect data.Value) {
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
