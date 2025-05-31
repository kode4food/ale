package compiler

import (
	"github.com/kode4food/ale/internal/compiler/encoder"
	"github.com/kode4food/ale/internal/types"
	"github.com/kode4food/ale/pkg/data"
)

// Call represents a code-generating function for the compiler
type Call func(encoder.Encoder, ...data.Value) error

var CallType = types.MakeBasic("special")

// Type makes Call a typed value
func (Call) Type() types.Type {
	return CallType
}

// Equal makes Call a typed Value
func (Call) Equal(data.Value) bool {
	return false
}

func (c Call) Get(key data.Value) (data.Value, bool) {
	return data.DumpMapped(c).Get(key)
}
