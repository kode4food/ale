package encoder

import "github.com/kode4food/ale/data"

// Call represents a code-generating function for the compiler
type Call func(Encoder, ...data.Value)

// Type makes Call a typed value
func (Call) Type() data.Name {
	return "encoder"
}

// Equal makes Call a typed Value
func (Call) Equal(_ data.Value) bool {
	return false
}

func (c Call) String() string {
	return data.DumpString(c)
}
