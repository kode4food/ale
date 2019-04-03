package test

import (
	"fmt"
	"testing"

	"gitlab.com/kode4food/ale/api"
	"gitlab.com/kode4food/ale/internal/compiler/arity"
)

func TestPredicates(t *testing.T) {
	testCode(t, `(eq true true true)`, api.True)
	testCode(t, `(eq true false true)`, api.False)
	testCode(t, `(eq false false false)`, api.True)

	testCode(t, `(!eq true true true)`, api.False)
	testCode(t, `(!eq true false)`, api.True)
	testCode(t, `(!eq false false)`, api.False)

	testCode(t, `(nil? nil)`, api.True)
	testCode(t, `(nil? nil nil nil)`, api.True)
	testCode(t, `(nil? () nil)`, api.False)
	testCode(t, `(nil? false)`, api.False)
	testCode(t, `(nil? false () nil)`, api.False)

	testCode(t, `(nil? "hello")`, api.False)
	testCode(t, `(nil? '(1 2 3))`, api.False)
	testCode(t, `(nil? () nil "hello")`, api.False)

	testCode(t, `(keyword? :hello)`, api.True)
	testCode(t, `(!keyword? :hello)`, api.False)
	testCode(t, `(keyword? 99)`, api.False)
	testCode(t, `(!keyword? 99)`, api.True)

	testBadCode(t, `(nil?)`, fmt.Errorf(arity.BadMinimumArity, 0, 1))
}
