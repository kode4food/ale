package macro

import (
	"gitlab.com/kode4food/ale/data"
	"gitlab.com/kode4food/ale/namespace"
)

// Call represents a macro's calling signature
type Call func(namespace.Type, ...data.Value) data.Value

// Type makes Call a typed value
func (Call) Type() data.Name {
	return "macro"
}

func (c Call) String() string {
	return data.DumpString(c)
}
