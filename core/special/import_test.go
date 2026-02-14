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

func TestImportErrors(t *testing.T) {
	as := assert.New(t)
	as.MustEval(`
		(define-namespace import-source
			(define x 99)
			(define y 100)
		)
	`)

	as.ErrorWith(`(import)`,
		fmt.Errorf(data.ErrRangedArity, 1, 2, 0),
	)
	as.ErrorWith(`(import import-source x y)`,
		fmt.Errorf(data.ErrRangedArity, 1, 2, 3),
	)
	as.ErrorWith(`(import 99)`,
		fmt.Errorf("%w: %s", special.ErrExpectedName, I(99)),
	)
	as.ErrorWith(`(import missing-source)`,
		fmt.Errorf(env.ErrNamespaceNotFound, LS("missing-source")),
	)
	as.ErrorWith(`(import import-source 99)`,
		fmt.Errorf("%w: %s", special.ErrUnexpectedImport, I(99)),
	)
	as.ErrorWith(`(import import-source (x x))`,
		fmt.Errorf("%w: %s", special.ErrDuplicateName, LS("x")),
	)
	as.ErrorWith(`(import import-source ([x]))`, special.ErrUnpairedBindings)
	as.ErrorWith(`(import import-source ([99 x]))`,
		fmt.Errorf("%w: %s", special.ErrExpectedName, I(99)),
	)
	as.ErrorWith(`(import import-source ([z 99]))`,
		fmt.Errorf("%w: %s", special.ErrExpectedName, I(99)),
	)
	as.ErrorWith(`(import import-source (99))`,
		fmt.Errorf("%w: %s", special.ErrUnexpectedImport, I(99)),
	)
}

func TestImportRuntimeFailures(t *testing.T) {
	as := assert.New(t)
	as.MustEval(`
		(define-namespace import-private
			(define :private hidden 42)
			(define shown 99)
		)
	`)

	as.PanicWith(`(import import-private hidden)`,
		fmt.Errorf(env.ErrNameNotDeclared, LS("hidden")),
	)
	as.PanicWith(`(import import-private missing)`,
		fmt.Errorf(env.ErrNameNotDeclared, LS("missing")),
	)

	as.MustEval(`
		(define-namespace import-collision
			(define value 1)
		)
	`)
	as.PanicWith(`
		(define value 99)
		(import import-collision value)
	`, fmt.Errorf(env.ErrNameAlreadyDeclared, data.Locals{"value"}))
}
