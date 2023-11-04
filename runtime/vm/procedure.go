package vm

import (
	"github.com/kode4food/ale/compiler/encoder"
	"github.com/kode4food/ale/compiler/ir/optimize"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/env"
	"github.com/kode4food/ale/internal/types"
	"github.com/kode4food/ale/runtime/isa"
)

// Procedure encapsulates the initial environment of an abstract machine
type Procedure struct {
	Globals      env.Namespace
	Constants    data.Values
	Code         isa.Instructions
	StackSize    int
	LocalCount   int
	ArityChecker data.ArityChecker
}

// MakeProcedure instantiates an abstract machine Procedure from the provided
// Encoder's intermediate representation
func MakeProcedure(e encoder.Encoder) *Procedure {
	code := e.Code()
	optimized := optimize.Instructions(code)
	return &Procedure{
		Globals:      e.Globals(),
		Constants:    e.Constants(),
		StackSize:    int(e.StackSize()),
		LocalCount:   int(e.LocalCount()),
		Code:         isa.Flatten(optimized),
		ArityChecker: data.AnyArityChecker,
	}
}

// Call allows an abstract machine Procedure to be called for the purpose of
// instantiating a closure. Only the compiler invokes this calling interface.
func (p *Procedure) Call(values ...data.Value) data.Value {
	return newClosure(p, values)
}

// CheckArity performs a compile-time arity check for the Procedure
func (p *Procedure) CheckArity(int) error {
	return nil
}

// Type makes Procedure a typed value
func (p *Procedure) Type() types.Type {
	return types.BasicProcedure
}

// Equal compares this Procedure to another for equality
func (p *Procedure) Equal(v data.Value) bool {
	return p == v
}

func (p *Procedure) String() string {
	return data.DumpString(p)
}
