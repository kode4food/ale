package special_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/pkg/core/special"
	"github.com/kode4food/ale/pkg/data"
)

func TestDefinitionsEval(t *testing.T) {
	as := assert.New(t)
	as.MustEvalTo(`
		(define foo "bar")
		foo
	`, S("bar"))

	as.MustEvalTo(`
		(define return-local
			(thunk (let [foo "local"] foo)))
		(return-local)
	`, S("local"))
}

func TestLetBindingErrors(t *testing.T) {
	as := assert.New(t)
	as.ErrorWith(`
		(let 99 "hello")
	`, fmt.Errorf(special.ErrUnexpectedLetSyntax, "99"))

	as.ErrorWith(`
		(let [a blah b] "hello")
	`, errors.New(special.ErrUnpairedBindings))

	as.PanicWith(`
		(let ((a blah)) "hello")
	`, unexpectedTypeError("list", "vector"))
}

func TestMutualBindingsEval(t *testing.T) {
	as := assert.New(t)

	as.MustEvalTo(`
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
