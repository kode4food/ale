package assert_test

import (
	"fmt"
	"testing"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/assert"
)

func TestEval(t *testing.T) {
	as := assert.New(t)
	as.EvalTo(`2`, data.Integer(2))
}

func TestPanicWith(t *testing.T) {
	as := assert.New(t)
	as.PanicWith(`(raise "boom!")`, fmt.Errorf("boom!"))
}
