package optimize_test

import (
	"testing"

	"github.com/kode4food/ale/compiler/encoder"
	"github.com/kode4food/ale/compiler/ir/optimize"
	"github.com/kode4food/ale/env"
	"github.com/kode4food/ale/internal/assert"
	"github.com/kode4food/ale/runtime/isa"
)

func TestLiterals(t *testing.T) {
	as := assert.New(t)

	ns := env.NewEnvironment().GetRoot()
	e1 := encoder.NewEncoder(ns)
	e1.Emit(isa.Null)
	e1.Emit(isa.Return)

	as.Instructions(isa.Instructions{
		isa.RetNull.New(),
	}, optimize.Encoded(e1.Encode()).Code)
}
