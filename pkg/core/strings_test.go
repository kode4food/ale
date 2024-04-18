package core_test

import (
	"testing"

	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/pkg/data"
)

func TestStrEval(t *testing.T) {
	as := assert.New(t)
	as.EvalTo(`
	  (str "hello" '() [1 2 3 4])
	`, S("hello[1 2 3 4]"))

	as.EvalTo(`
	  (string? "hello" "there")
	`, data.True)

	as.EvalTo(`
	  (string? "hello" 99)
	`, data.False)
}

func TestReadableStrEval(t *testing.T) {
	as := assert.New(t)
	as.EvalTo("(str! \"hello\nyou\")", S("\"hello\\nyou\""))
	as.EvalTo(`(str! "hello" "you")`, S(`"hello" "you"`))
	as.EvalTo(`(str!)`, S(""))
}
