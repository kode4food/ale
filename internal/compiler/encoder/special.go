package encoder

import "gitlab.com/kode4food/ale/api"

// Call represents a code-generating function for the Compiler
type Call func(Type, ...api.Value)

// Type makes Special a typed value
func (Call) Type() api.Name {
	return "#special"
}

func (c Call) String() string {
	return api.DumpString(c)
}
