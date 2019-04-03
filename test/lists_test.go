package test

import (
	"testing"

	"gitlab.com/kode4food/ale/api"
	. "gitlab.com/kode4food/ale/internal/assert/helpers"
)

func TestList(t *testing.T) {
	testCode(t, `(list? '(1 2 3))`, api.True)
	testCode(t, `(list? ())`, api.True)
	testCode(t, `(list? [1 2 3])`, api.False)
	testCode(t, `(list? 42)`, api.False)
	testCode(t, `(list? (list 1 2 3))`, api.True)
	testCode(t, `(list)`, api.EmptyList)

	testCode(t, `
		(def x '(1 2 3 4))
		(x 2)
	`, F(3))
}
