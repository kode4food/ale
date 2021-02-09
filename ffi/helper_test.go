package ffi_test

import (
	"testing"

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

	Env map[data.Name]interface{}
)

var (
	testEnv = env.NewEnvironment()
	ready   bool
)

func NewWrapped(t *testing.T) *EvalWrapped {
	return &EvalWrapped{
		Wrapper: assert.New(t),
	}
}

func (e *EvalWrapped) EvalTo(src string, env Env, expect data.Value) {
	e.Helper()
	if !ready {
		bootstrap.Into(testEnv)
		ready = true
	}
	ns := testEnv.GetAnonymous()
	for n, v := range env {
		v, err := ffi.Wrap(v)
		e.Nil(err)
		ns.Declare(n).Bind(v)
	}
	res := eval.String(ns, data.String(src))
	e.Equal(expect, res)
}
