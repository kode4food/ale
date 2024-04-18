package procedure

import (
	"github.com/kode4food/ale/pkg/compiler/encoder"
	"github.com/kode4food/ale/pkg/compiler/ir/analysis"
	"github.com/kode4food/ale/pkg/compiler/ir/optimize"
	"github.com/kode4food/ale/pkg/data"
	"github.com/kode4food/ale/pkg/runtime/vm"
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
