package builtin_test

import (
	"testing"

	"gitlab.com/kode4food/ale/data"
	"gitlab.com/kode4food/ale/internal/assert"
	. "gitlab.com/kode4food/ale/internal/assert/helpers"
)

func TestIfEval(t *testing.T) {
	as := assert.New(t)
	as.EvalTo(`(if #f 1 0)`, F(0))
	as.EvalTo(`(if #t 1 0)`, F(1))
	as.EvalTo(`(if '() 1 0)`, F(0))
	as.EvalTo(`(if "hello" 1 0)`, F(1))
	as.EvalTo(`(if #f 1)`, data.Null)
}
