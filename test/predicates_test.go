package test

import (
	"fmt"
	"testing"

	"gitlab.com/kode4food/ale/compiler/arity"
	"gitlab.com/kode4food/ale/data"
)

func TestPredicates(t *testing.T) {
	testCode(t, `(eq true true true)`, data.True)
	testCode(t, `(eq true false true)`, data.False)
	testCode(t, `(eq false false false)`, data.True)

	testCode(t, `(!eq true true true)`, data.False)
	testCode(t, `(!eq true false)`, data.True)
	testCode(t, `(!eq false false)`, data.False)

	testCode(t, `(nil? nil)`, data.True)
	testCode(t, `(nil? nil nil nil)`, data.True)
	testCode(t, `(nil? () nil)`, data.False)
	testCode(t, `(nil? false)`, data.False)
	testCode(t, `(nil? false () nil)`, data.False)

	testCode(t, `(nil? "hello")`, data.False)
	testCode(t, `(nil? '(1 2 3))`, data.False)
	testCode(t, `(nil? () nil "hello")`, data.False)

	testCode(t, `(keyword? :hello)`, data.True)
	testCode(t, `(!keyword? :hello)`, data.False)
	testCode(t, `(keyword? 99)`, data.False)
	testCode(t, `(!keyword? 99)`, data.True)

	testBadCode(t, `(nil?)`, fmt.Errorf(arity.BadMinimumArity, 0, 1))
}
