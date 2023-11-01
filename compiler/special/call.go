package special

import (
	"github.com/kode4food/ale/compiler/encoder"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/types"
)

// Call represents a code-generating function for the compiler
type Call func(encoder.Encoder, ...data.Value)

var CallType = types.MakeBasic("special")

// Type makes Call a typed value
func (Call) Type() types.Type {
	return CallType
}

// Equal makes Call a typed Value
func (Call) Equal(data.Value) bool {
	return false
}

func (c Call) String() string {
	return data.DumpString(c)
}
