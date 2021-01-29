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

// Equal compares this Call to another for equality
func (Call) Equal(_ data.Value) bool {
	return false
}

func (c Call) String() string {
	return data.DumpString(c)
}
