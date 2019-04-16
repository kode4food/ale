package encoder

import "gitlab.com/kode4food/ale/api"

// Call represents a code-generating function for the compiler
type Call func(Type, ...api.Value)

// Type makes Call a typed value
func (Call) Type() api.Name {
	return "Special"
}

func (c Call) String() string {
	return api.DumpString(c)
}
