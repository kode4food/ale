package vm

import (
	"github.com/kode4food/ale/compiler/encoder"
	"github.com/kode4food/ale/compiler/ir/optimize"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/env"
	"github.com/kode4food/ale/runtime/isa"
	"github.com/kode4food/ale/types"
)

// Lambda encapsulates the initial environment of an abstract machine
type Lambda struct {
	Globals      env.Namespace
	Constants    data.Values
	Code         isa.Instructions
	StackSize    int
	LocalCount   int
	ArityChecker data.ArityChecker
}

// LambdaFromEncoder instantiates an abstract machine Lambda from the provided
// Encoder's intermediate representation
func LambdaFromEncoder(e encoder.Encoder) *Lambda {
	code := e.Code()
	optimized := optimize.Instructions(code)
	return &Lambda{
		Globals:      e.Globals(),
		Constants:    e.Constants(),
		StackSize:    int(e.StackSize()),
		LocalCount:   int(e.LocalCount()),
		Code:         isa.Flatten(optimized),
		ArityChecker: NoArityChecker,
	}
}

func NoArityChecker(int) error {
	return nil
}

// Call allows an abstract machine Lambda to be called for the purpose of
// instantiating a closure. Only the compiler invokes this calling interface.
func (l *Lambda) Call(values ...data.Value) data.Value {
	return newClosure(l, values)
}

// CheckArity performs a compile-time arity check for the Function
func (l *Lambda) CheckArity(int) error {
	return nil
}

// Convention returns the Function's calling convention
func (l *Lambda) Convention() data.Convention {
	return data.NormalCall
}

// Type makes Lambda a typed value
func (l *Lambda) Type() types.Type {
	return types.Lambda
}

// Equal compares this Lambda to another for equality
func (l *Lambda) Equal(v data.Value) bool {
	return l == v
}

func (l *Lambda) String() string {
	return data.DumpString(l)
}
