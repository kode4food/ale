package builtin_test

import (
	"testing"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
)

func TestStrEval(t *testing.T) {
	as := assert.New(t)
	as.MustEvalTo(`
	  (str "hello" '() [1 2 3 4])
	`, S("hello[1 2 3 4]"))

	as.MustEvalTo(`
	  (string? "hello" "there")
	`, data.True)

	as.MustEvalTo(`
	  (string? "hello" 99)
	`, data.False)
}

func TestReadableStrEval(t *testing.T) {
	as := assert.New(t)
	as.MustEvalTo("(str! \"hello\nyou\")", S("\"hello\\nyou\""))
	as.MustEvalTo(`(str! "hello" "you")`, S(`"hello" "you"`))
	as.MustEvalTo(`(str!)`, S(""))
}
