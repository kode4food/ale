package test

import (
	"testing"

	"gitlab.com/kode4food/ale/data"
	. "gitlab.com/kode4food/ale/internal/assert/helpers"
)

func TestList(t *testing.T) {
	testCode(t, `(list? '(1 2 3))`, data.True)
	testCode(t, `(list? ())`, data.True)
	testCode(t, `(list? [1 2 3])`, data.False)
	testCode(t, `(list? 42)`, data.False)
	testCode(t, `(list? (list 1 2 3))`, data.True)
	testCode(t, `(list)`, data.EmptyList)

	testCode(t, `
		(def x '(1 2 3 4))
		(x 2)
	`, F(3))
}
