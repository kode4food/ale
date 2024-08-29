package optimize_test

import (
	"testing"

	"github.com/kode4food/ale/internal/assert"
	"github.com/kode4food/ale/internal/compiler/encoder"
	"github.com/kode4food/ale/internal/compiler/ir/optimize"
	"github.com/kode4food/ale/internal/runtime/isa"
	"github.com/kode4food/ale/pkg/env"
)

func TestIneffectivePushes(t *testing.T) {
	as := assert.New(t)

	ns := env.NewEnvironment().GetRoot()
	e1 := encoder.NewEncoder(ns)
	e1.Emit(isa.True)
	e1.Emit(isa.Const, 0)
	e1.Emit(isa.Pop)
	e1.Emit(isa.Return)

	as.Instructions(isa.Instructions{
		isa.RetTrue.New(),
	}, optimize.Encoded(e1.Encode()).Code)
}
