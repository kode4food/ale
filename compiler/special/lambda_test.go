package special_test

import (
	"testing"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/assert"
)

func TestEmptyLambda(t *testing.T) {
	as := assert.New(t)

	as.EvalTo(`((lambda))`, data.Null)
}
