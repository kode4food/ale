package optimize_test

import (
	"testing"

	"github.com/kode4food/ale/env"
	"github.com/kode4food/ale/internal/assert"
	"github.com/kode4food/ale/internal/compiler/encoder"
	"github.com/kode4food/ale/internal/compiler/ir/optimize"
	"github.com/kode4food/ale/internal/runtime/isa"
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
