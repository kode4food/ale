package vm

import (
	"github.com/kode4food/ale/compiler/encoder"
	"github.com/kode4food/ale/compiler/ir/optimize"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/namespace"
	"github.com/kode4food/ale/runtime/isa"
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

// Call allows a VM Lambda to be called for the purpose
// of instantiating a closure. This calling interface is used
// only by the compiler.
func (l *Lambda) Call() data.Call {
	return func(values ...data.Value) data.Value {
		return newClosure(l, values)
	}
}

func (l *Lambda) String() string {
	return data.DumpString(l)
}
