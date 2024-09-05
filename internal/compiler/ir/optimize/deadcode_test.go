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

func TestIneffectiveStores(t *testing.T) {
	as := assert.New(t)

	ns := env.NewEnvironment().GetRoot()
	e1 := encoder.NewEncoder(ns)
	e1.Emit(isa.True)
	e1.Emit(isa.Store, 0)
	e1.Emit(isa.Load, 0)
	e1.Emit(isa.Return)

	as.Instructions(isa.Instructions{
		isa.RetTrue.New(),
	}, optimize.Encoded(e1.Encode()).Code)

	e2 := encoder.NewEncoder(ns)
	e2.Emit(isa.PosInt, 1)
	e2.Emit(isa.Store, 0)
	e2.Emit(isa.Load, 0)
	e2.Emit(isa.Load, 0)
	e2.Emit(isa.Add)
	e2.Emit(isa.Store, 1)
	e2.Emit(isa.Load, 1)
	e2.Emit(isa.Return)

	as.Instructions(isa.Instructions{
		isa.PosInt.New(1),
		isa.Store.New(0),
		isa.Load.New(0),
		isa.Load.New(0),
		isa.Add.New(),
		isa.Return.New(),
	}, optimize.Encoded(e2.Encode()).Code)

	e3 := encoder.NewEncoder(ns)
	e3.Emit(isa.PosInt, 1)
	e3.Emit(isa.PosInt, 2)
	e3.Emit(isa.Store, 0)
	e3.Emit(isa.Store, 1)
	e3.Emit(isa.Load, 1)
	e3.Emit(isa.Load, 0)
	e3.Emit(isa.Add)
	e3.Emit(isa.Return)

	as.Instructions(isa.Instructions{
		isa.PosInt.New(1),
		isa.PosInt.New(2),
		isa.Add.New(),
		isa.Return.New(),
	}, optimize.Encoded(e3.Encode()).Code)
}
