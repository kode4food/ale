package special_test

import (
	"testing"

	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/pkg/data"
)

func TestIfEval(t *testing.T) {
	as := assert.New(t)
	as.MustEvalTo(`(if false 1 0)`, F(0))
	as.MustEvalTo(`(if true 1 0)`, F(1))
	as.MustEvalTo(`(if '() 1 0)`, F(1))
	as.MustEvalTo(`(if "hello" 1 0)`, F(1))
	as.MustEvalTo(`(if false 1)`, data.Null)
}
