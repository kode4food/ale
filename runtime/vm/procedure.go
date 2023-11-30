package vm

import (
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/env"
	"github.com/kode4food/ale/internal/types"
	"github.com/kode4food/ale/runtime/isa"
)

type (
	// Procedure encapsulates the initial environment of an abstract machine
	Procedure struct {
		Globals      env.Namespace
		Constants    data.Vector
		Code         isa.Instructions
		StackSize    int
		LocalCount   int
		ArityChecker data.ArityChecker
	}

	Closure struct {
		*Procedure
		Captured data.Vector
	}
)

// Call allows an abstract machine Procedure to be called for the purpose of
// instantiating a Closure. Only the compiler invokes this calling interface.
func (p *Procedure) Call(values ...data.Value) data.Value {
	return &Closure{
		Procedure: p,
		Captured:  values,
	}
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

func (p *Procedure) Get(key data.Value) (data.Value, bool) {
	return data.DumpMapped(p).Get(key)
}

// Call turns Closure into a Procedure, and serves as the virtual machine
func (c *Closure) Call(args ...data.Value) data.Value {
	return (&VM{CL: c, ARGS: args}).Run()
}

// CheckArity performs a compile-time arity check for the Closure
func (c *Closure) CheckArity(i int) error {
	return c.ArityChecker(i)
}

func (c *Closure) Equal(v data.Value) bool {
	return c == v
}
