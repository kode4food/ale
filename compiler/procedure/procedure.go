package procedure

import (
	"github.com/kode4food/ale/compiler/encoder"
	"github.com/kode4food/ale/compiler/ir/optimize"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/runtime/vm"
)

// FromEncoded instantiates an abstract machine Procedure from the provided
// Encoded representation
func FromEncoded(e *encoder.Encoded) *vm.Procedure {
	return &vm.Procedure{
		Runnable:     *optimize.Encoded(e).Runnable(),
		ArityChecker: data.AnyArityChecker,
	}
}
