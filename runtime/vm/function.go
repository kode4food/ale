package vm

import (
	"gitlab.com/kode4food/ale/compiler/encoder"
	"gitlab.com/kode4food/ale/compiler/ir/optimize"
	"gitlab.com/kode4food/ale/data"
	"gitlab.com/kode4food/ale/namespace"
	"gitlab.com/kode4food/ale/runtime/isa"
)

// Function encapsulates the initial environment of a virtual machine
type Function struct {
	Globals    namespace.Type
	Constants  data.Values
	Code       []isa.Word
	StackSize  int
	LocalCount int
}

// FunctionFromEncoder instantiates a VM Function from the provided
// Encoder's intermediate representation
func FunctionFromEncoder(e encoder.Type) *Function {
	code := e.Code()
	optimized := optimize.Instructions(code)
	return &Function{
		Globals:    e.Globals(),
		Constants:  e.Constants(),
		StackSize:  e.StackSize(),
		LocalCount: e.LocalCount(),
		Code:       isa.Flatten(optimized),
	}
}

// Caller allows a VM Function to be called for the purpose
// of instantiating a Closure. This calling interface is used
// only by the compiler.
func (f *Function) Caller() data.Call {
	return func(values ...data.Value) data.Value {
		closure := &Closure{
			Function: f,
			Values:   values,
		}
		return closure.Caller()
	}
}

func (f *Function) String() string {
	return data.DumpString(f)
}
