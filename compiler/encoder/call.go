package encoder

import "github.com/kode4food/ale/data"

// Call represents a code-generating function for the compiler
type Call func(Type, ...data.Value)

// Type makes Call a typed value
func (Call) Type() data.Name {
	return "encoder-call"
}

func (c Call) String() string {
	return data.DumpString(c)
}
