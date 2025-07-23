package special_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/kode4food/ale/core/special"
	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
)

func TestLet(t *testing.T) {
	as := assert.New(t)
	as.MustEvalTo(`(let ([x 99][y 1]) (+ x y))`, I(100))
}

func TestBindingErrors(t *testing.T) {
	as := assert.New(t)

	as.ErrorWith(`(let ([x 99][y]) (+ x y))`,
		errors.New(special.ErrUnpairedBindings),
	)

	as.ErrorWith(`(let ([x 99][x 99]) (+ x y))`,
		fmt.Errorf(special.ErrNameAlreadyBound, "x"),
	)

	as.ErrorWith(`(let (x . 99) (+ x x))`,
		fmt.Errorf(special.ErrUnexpectedLetSyntax, "(x . 99)"),
	)
}
