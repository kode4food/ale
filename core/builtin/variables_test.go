package builtin_test

import (
	"fmt"
	"testing"

	"gitlab.com/kode4food/ale/compiler/special"
	"gitlab.com/kode4food/ale/data"
	"gitlab.com/kode4food/ale/internal/assert"
	. "gitlab.com/kode4food/ale/internal/assert/helpers"
)

func TestDefinitionsEval(t *testing.T) {
	as := assert.New(t)
	as.EvalTo(`
		(define foo "bar")
		foo
	`, S("bar"))

	as.EvalTo(`
		(define return-local (lambda ()
			(let [foo "local"] foo)))
		(return-local)
	`, S("local"))
}

func TestLetBindingErrors(t *testing.T) {
	as := assert.New(t)
	as.PanicWith(`
		(let 99 "hello")
	`, fmt.Errorf(special.UnexpectedLetSyntax, "99"))

	as.PanicWith(`
		(let [a blah b] "hello")
	`, fmt.Errorf(special.UnpairedBindings))

	as.PanicWith(`
		(let ((a blah)) "hello")
	`, typeErr("*data.list", "data.Vector"))
}

func TestMutualBindingsEval(t *testing.T) {
	as := assert.New(t)

	as.EvalTo(`
		(letrec ([
			is-even?
			(lambda (n) (or (= n 0)
			                (is-odd? (dec n))))]

			[is-odd?
			(lambda (n) (and (not (= n 0))
			                 (is-even? (dec n))))])
		(is-even? 13))
	`, data.False)
}
