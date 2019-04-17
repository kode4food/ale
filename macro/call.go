package macro

import (
	"gitlab.com/kode4food/ale/api"
	"gitlab.com/kode4food/ale/namespace"
)

// Call represents a macro's calling signature
type Call func(namespace.Type, ...api.Value) api.Value

// Type makes Call a typed value
func (Call) Type() api.Name {
	return "Macro"
}

func (c Call) String() string {
	return api.DumpString(c)
}
