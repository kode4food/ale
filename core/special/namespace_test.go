package special_test

import (
	"fmt"
	"testing"

	"github.com/kode4food/ale/core/special"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/env"
	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
)

func TestNamespaceDefinition(t *testing.T) {
	as := assert.New(t)

	d := as.MustEval(`
		(define-namespace some-namespace
			(define x 99)
			(define y 100)
		)
	`).(data.Vector)
	_, ok1 := d.IndexOf(LS("x"))
	_, ok2 := d.IndexOf(LS("y"))
	as.True(ok1 && ok2)

	i := as.MustEval(`(import some-namespace)`).(data.Vector)
	_, ok1 = i.IndexOf(LS("x"))
	_, ok2 = i.IndexOf(LS("y"))
	as.True(ok1 && ok2)

	as.MustEvalTo(`(import some-namespace y) y`, I(100))
	as.MustEvalTo(`(import some-namespace [x99 x]) x99`, I(99))
	as.MustEvalTo(`(import some-namespace (x y)) [y x]`, V(I(100), I(99)))
}

func TestDeclaredCurrentNamespace(t *testing.T) {
	as := assert.New(t)
	d := as.MustEval(`
		(define declared-probe 42)
		(declared)
	`).(data.Vector)
	_, ok := d.IndexOf(LS("declared-probe"))
	as.True(ok)
}

func TestDeclaredErrors(t *testing.T) {
	as := assert.New(t)
	as.ErrorWith(`(declared some-namespace also-nope)`,
		fmt.Errorf(data.ErrRangedArity, 0, 1, 2),
	)
	as.ErrorWith(`(declared 99)`,
		fmt.Errorf("%w: %s", special.ErrExpectedName, I(99)),
	)
	as.ErrorWith(`(declared does-not-exist)`,
		fmt.Errorf(env.ErrNamespaceNotFound, LS("does-not-exist")),
	)
}

func TestMakeNamespaceErrors(t *testing.T) {
	as := assert.New(t)

	as.ErrorWith(`(%mk-ns)`,
		fmt.Errorf(data.ErrFixedArity, 1, 0),
	)
	as.ErrorWith(`(%mk-ns 99)`,
		fmt.Errorf("%w: %s", special.ErrExpectedName, I(99)),
	)

	as.MustEval(`(%mk-ns duplicate-ns)`)
	as.PanicWith(`(%mk-ns duplicate-ns)`,
		fmt.Errorf(env.ErrNamespaceExists, LS("duplicate-ns")),
	)

	as.MustEval(`(define in-ns (%mk-ns bad-form-ns))`)
	as.PanicWith(`(in-ns unknown-name)`,
		fmt.Errorf(env.ErrNameNotDeclared, LS("unknown-name")),
	)
}
