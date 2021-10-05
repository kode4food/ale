package encoder

import (
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/types"
	"github.com/kode4food/ale/types/basic"
)

// Call represents a code-generating function for the compiler
type Call func(Encoder, ...data.Value)

var encoderType = basic.New("encoder")

// Type makes Call a typed value
func (Call) Type() types.Type {
	return encoderType
}

// Equal makes Call a typed Value
func (Call) Equal(_ data.Value) bool {
	return false
}

func (c Call) String() string {
	return data.DumpString(c)
}
