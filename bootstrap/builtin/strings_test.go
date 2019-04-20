package builtin_test

import (
	"testing"

	"gitlab.com/kode4food/ale/data"
	"gitlab.com/kode4food/ale/internal/assert"
	. "gitlab.com/kode4food/ale/internal/assert/helpers"
)

func TestStrEval(t *testing.T) {
	as := assert.New(t)
	as.EvalTo(`
	  (str "hello" nil [1 2 3 4])
	`, S("hello[1 2 3 4]"))

	as.EvalTo(`
	  (str? "hello" "there")
	`, data.True)

	as.EvalTo(`
	  (str? "hello" 99)
	`, data.False)
}

func TestReadableStrEval(t *testing.T) {
	as := assert.New(t)
	as.EvalTo("(str! \"hello\nyou\")", S("\"hello\\nyou\""))
	as.EvalTo(`(str! "hello" "you")`, S(`"hello" "you"`))
}
