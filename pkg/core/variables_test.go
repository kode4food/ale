package core_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	builtin "github.com/kode4food/ale/pkg/core"
	"github.com/kode4food/ale/pkg/data"
)

func TestDefinitionsEval(t *testing.T) {
	as := assert.New(t)
	as.EvalTo(`
		(define foo "bar")
		foo
	`, S("bar"))

	as.EvalTo(`
		(define return-local
			(thunk (let [foo "local"] foo)))
		(return-local)
	`, S("local"))
}

func TestLetBindingErrors(t *testing.T) {
	as := assert.New(t)
	as.PanicWith(`
		(let 99 "hello")
	`, fmt.Errorf(builtin.ErrUnexpectedLetSyntax, "99"))

	as.PanicWith(`
		(let [a blah b] "hello")
	`, errors.New(builtin.ErrUnpairedBindings))

	as.PanicWith(`
		(let ((a blah)) "hello")
	`, unexpectedTypeError("list", "vector"))
}

func TestMutualBindingsEval(t *testing.T) {
	as := assert.New(t)

	as.EvalTo(`
		(let-rec ([
			is-even?
			(lambda (n) (or (= n 0)
			                (is-odd? (dec n))))]

			[is-odd?
			(lambda (n) (and (not (= n 0))
			                 (is-even? (dec n))))])
		(is-even? 13))
	`, data.False)
}
