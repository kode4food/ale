package test

import (
	"testing"

	"gitlab.com/kode4food/ale/data"
	. "gitlab.com/kode4food/ale/internal/assert/helpers"
)

func TestGenerate(t *testing.T) {
	testCode(t, `
		(def g (generate
			(emit 99)
			(emit 100 1000)))
		(apply + g)
	`, F(1199))
}

func TestPromise(t *testing.T) {
	testCode(t, `
		(def p1 (promise))
		(promise? p1)
	`, data.True)

	testCode(t, `
		(def p2 (promise "hello"))
		(p2)
	`, S("hello"))
}

func TestFuture(t *testing.T) {
	testCode(t, `
		(def p (future "hello"))
		(p)
	`, S("hello"))
}
