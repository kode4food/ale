package vm

import (
	"gitlab.com/kode4food/ale/compiler/encoder"
	"gitlab.com/kode4food/ale/compiler/ir/optimize"
	"gitlab.com/kode4food/ale/data"
	"gitlab.com/kode4food/ale/namespace"
	"gitlab.com/kode4food/ale/runtime/isa"
)

// Config encapsulates the initial environment of a virtual machine
type Config struct {
	Globals    namespace.Type
	Constants  data.Values
	Code       []isa.Word
	StackSize  int
	LocalCount int
}

// ConfigFromEncoder optimizes and flattens a VM config from the
// provided encoder's intermediate representation
func ConfigFromEncoder(e encoder.Type) *Config {
	code := e.Code()
	optimized := optimize.Instructions(code)
	return &Config{
		Globals:    e.Globals(),
		Constants:  e.Constants(),
		StackSize:  e.StackSize(),
		LocalCount: e.LocalCount(),
		Code:       isa.Flatten(optimized),
	}
}
