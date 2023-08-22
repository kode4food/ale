package encoder

import (
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/types"
)

// Call represents a code-generating function for the compiler
type Call func(Encoder, ...data.Value)

var encoderType = types.Basic("encoder")

// Type makes Call a typed value
func (Call) Type() types.Type {
	return encoderType
}

// Equal makes Call a typed Value
func (Call) Equal(data.Value) bool {
	return false
}

func (c Call) String() string {
	return data.DumpString(c)
}
