package special_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/kode4food/ale/compiler/special"
	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
)

func testEvalError(t *testing.T, src string, err error) {
	as := assert.New(t)
	defer as.ExpectPanic(err)
	as.Eval(src)
}

func TestLet(t *testing.T) {
	as := assert.New(t)
	as.EvalTo(`(let ([x 99][y 1]) (+ x y))`, I(100))
}

func TestBindingErrors(t *testing.T) {
	testEvalError(t, `(let ([x 99][y]) (+ x y))`,
		errors.New(special.ErrUnpairedBindings),
	)

	testEvalError(t, `(let ([x 99][x 99]) (+ x y))`,
		fmt.Errorf(special.ErrNameAlreadyBound, "x"),
	)

	testEvalError(t, `(let (x . 99) (+ x x))`,
		fmt.Errorf(special.ErrUnexpectedLetSyntax, "(x . 99)"),
	)
}
