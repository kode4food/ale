package test

import (
	"testing"

	"gitlab.com/kode4food/ale/api"
	. "gitlab.com/kode4food/ale/internal/assert/helpers"
)

func TestIf(t *testing.T) {
	testCode(t, `(if false 1 0)`, F(0))
	testCode(t, `(if true 1 0)`, F(1))
	testCode(t, `(if nil 1 0)`, F(0))
	testCode(t, `(if () 1 0)`, F(1))
	testCode(t, `(if "hello" 1 0)`, F(1))
	testCode(t, `(if false 1)`, api.Nil)
}
