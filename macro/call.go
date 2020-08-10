package macro

import (
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/env"
)

// Call represents a macro's calling signature
type Call func(env.Namespace, ...data.Value) data.Value

// Type makes Call a typed value
func (Call) Type() data.Name {
	return "macro"
}

func (c Call) String() string {
	return data.DumpString(c)
}
