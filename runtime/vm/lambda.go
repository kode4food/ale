package vm

import (
	"gitlab.com/kode4food/ale/compiler/encoder"
	"gitlab.com/kode4food/ale/compiler/ir/optimize"
	"gitlab.com/kode4food/ale/data"
	"gitlab.com/kode4food/ale/namespace"
	"gitlab.com/kode4food/ale/runtime/isa"
)

// Lambda encapsulates the initial environment of a virtual machine
type Lambda struct {
	Globals      namespace.Type
	Constants    data.Values
	Code         []isa.Word
	StackSize    int
	LocalCount   int
	ArityChecker data.ArityChecker
}

// LambdaFromEncoder instantiates a VM Lambda from the provided
// Encoder's intermediate representation
func LambdaFromEncoder(e encoder.Type) *Lambda {
	code := e.Code()
	optimized := optimize.Instructions(code)
	return &Lambda{
		Globals:    e.Globals(),
		Constants:  e.Constants(),
		StackSize:  e.StackSize(),
		LocalCount: e.LocalCount(),
		Code:       isa.Flatten(optimized),
	}
}

// Caller allows a VM Lambda to be called for the purpose
// of instantiating a Closure. This calling interface is used
// only by the compiler.
func (l *Lambda) Caller() data.Call {
	return func(values ...data.Value) data.Value {
		return newClosure(l, values)
	}
}

func (l *Lambda) String() string {
	return data.DumpString(l)
}
