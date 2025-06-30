package special_test

import (
	"testing"

	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/pkg/data"
)

func TestNamespaceDefinition(t *testing.T) {
	as := assert.New(t)

	as.MustEvalTo(`
		(define-namespace some-namespace
			(define x 99)
			(define y 100)
		)
		(import some-namespace x)
		x
	`, I(99))

	r := as.MustEval(`(import some-namespace)`).(data.Vector)
	_, ok := r.IndexOf(LS("x"))
	as.True(ok)
	_, ok = r.IndexOf(LS("y"))
	as.True(ok)

	as.MustEvalTo(`(import some-namespace y) y`, I(100))
	as.MustEvalTo(`(import some-namespace [x99 x]) x99`, I(99))
}
