package core_test

import (
	"fmt"
	"testing"

	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/pkg/core"
	"github.com/kode4food/ale/pkg/data"
)

func TestLambda(t *testing.T) {
	as := assert.New(t)

	as.EvalTo(`((lambda [(x y) (+ x y)] [(x) (* x 2)]) 20 30)`, I(50))
	as.EvalTo(`((lambda [(x y) (+ x y)] [(x) (* x 2)]) 20)`, I(40))

	as.EvalTo(
		`((lambda [(x y) (+ x y)] [(x) (* x 2)] [x x]) 20 30 40)`,
		V(I(20), I(30), I(40)),
	)

	as.EvalTo(`((lambda x x) 1 2 3)`, V(I(1), I(2), I(3)))
	as.EvalTo(`((lambda (x) x) 1)`, I(1))
	as.EvalTo(`((lambda))`, data.Null)
}

func TestLambdaErrors(t *testing.T) {
	as := assert.New(t)

	as.PanicWith(`(lambda :kwd '())`,
		fmt.Errorf(core.ErrUnexpectedCaseSyntax, ":kwd"),
	)

	as.PanicWith(`(lambda [:kwd] '())`,
		fmt.Errorf(core.ErrUnexpectedParamSyntax, ":kwd"),
	)
}
