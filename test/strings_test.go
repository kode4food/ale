package test

import (
	"testing"

	"gitlab.com/kode4food/ale/api"
	. "gitlab.com/kode4food/ale/internal/assert/helpers"
)

func TestStr(t *testing.T) {
	testCode(t, `
	  (str "hello" nil [1 2 3 4])
	`, S("hello[1 2 3 4]"))

	testCode(t, `
	  (str? "hello" "there")
	`, api.True)

	testCode(t, `
	  (str? "hello" 99)
	`, api.False)
}

func TestReadableStr(t *testing.T) {
	testCode(t, "(str! \"hello\nyou\")", S("\"hello\nyou\""))
	testCode(t, `(str! "hello" "you")`, S(`"hello" "you"`))
}
