package procedure

import (
	"github.com/kode4food/ale/compiler/encoder"
	"github.com/kode4food/ale/compiler/ir/analysis"
	"github.com/kode4food/ale/compiler/ir/optimize"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/runtime/vm"
)

// FromEncoded instantiates an abstract machine Procedure from the provided
// Encoded representation
func FromEncoded(e *encoder.Encoded) *vm.Procedure {
	analysis.MustVerify(e.Code)
	return vm.MakeProcedure(
		optimize.Encoded(e).Runnable(),
		data.AnyArityChecker,
	)
}
