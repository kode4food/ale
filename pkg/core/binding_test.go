package core_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/pkg/core"
)

func TestLet(t *testing.T) {
	as := assert.New(t)
	as.EvalTo(`(let ([x 99][y 1]) (+ x y))`, I(100))
}

func TestBindingErrors(t *testing.T) {
	as := assert.New(t)

	as.PanicWith(`(let ([x 99][y]) (+ x y))`,
		errors.New(core.ErrUnpairedBindings),
	)

	as.PanicWith(`(let ([x 99][x 99]) (+ x y))`,
		fmt.Errorf(core.ErrNameAlreadyBound, "x"),
	)

	as.PanicWith(`(let (x . 99) (+ x x))`,
		fmt.Errorf(core.ErrUnexpectedLetSyntax, "(x . 99)"),
	)
}
