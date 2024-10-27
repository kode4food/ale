package procedure

import (
	"github.com/kode4food/ale/internal/compiler/encoder"
	"github.com/kode4food/ale/internal/compiler/ir/analysis"
	"github.com/kode4food/ale/internal/compiler/ir/optimize"
	"github.com/kode4food/ale/internal/runtime/vm"
	"github.com/kode4food/ale/pkg/data"
)

// FromEncoded instantiates an abstract machine Procedure from the provided
// Encoded representation
func FromEncoded(e *encoder.Encoded) (*vm.Procedure, error) {
	if err := analysis.Verify(e.Code); err != nil {
		return nil, err
	}
	run, err := optimize.Encoded(e).Runnable()
	if err != nil {
		return nil, err
	}
	return vm.MakeProcedure(run, data.CheckAnyArity), nil
}
