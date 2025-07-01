package special_test

import (
	"testing"

	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/pkg/data"
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
