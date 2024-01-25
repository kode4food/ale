package procedure

import (
	"github.com/kode4food/ale/compiler/encoder"
	"github.com/kode4food/ale/compiler/ir/optimize"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/runtime/isa"
	"github.com/kode4food/ale/runtime/vm"
)

// FromEncoded instantiates an abstract machine Procedure from the provided
// Encoded representation
func FromEncoded(e *encoder.Encoded) *vm.Procedure {
	optimized := optimize.FromEncoded(e)
	return &vm.Procedure{
		Globals:      e.Globals,
		Constants:    e.Constants,
		StackSize:    int(e.StackSize),
		LocalCount:   int(e.LocalCount),
		Code:         isa.Flatten(optimized),
		ArityChecker: data.AnyArityChecker,
	}
}
