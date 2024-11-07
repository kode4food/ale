package special_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/pkg/core/internal"
)

func TestLambda(t *testing.T) {
	as := assert.New(t)

	as.MustEvalTo(`((lambda [(x y) (+ x y)] [(x) (* x 2)]) 20 30)`, I(50))
	as.MustEvalTo(`((lambda [(x y) (+ x y)] [(x) (* x 2)]) 20)`, I(40))

	as.MustEvalTo(
		`((lambda [(x y) (+ x y)] [(x) (* x 2)] [x x]) 20 30 40)`,
		V(I(20), I(30), I(40)),
	)

	as.MustEvalTo(`((lambda x x) 1 2 3)`, V(I(1), I(2), I(3)))
	as.MustEvalTo(`((lambda (x) x) 1)`, I(1))
}

func TestLambdaErrors(t *testing.T) {
	as := assert.New(t)

	as.ErrorWith(`(lambda :kwd '())`,
		fmt.Errorf(internal.ErrUnexpectedCaseSyntax, ":kwd"),
	)

	as.ErrorWith(`(lambda [:kwd] '())`,
		fmt.Errorf(internal.ErrUnexpectedParamSyntax, ":kwd"),
	)

	as.ErrorWith(`(lambda)`,
		errors.New(internal.ErrNoCasesDefined),
	)
}
