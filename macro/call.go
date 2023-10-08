package macro

import (
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/env"
	"github.com/kode4food/ale/internal/types"
)

// Call represents a macro's calling signature
type Call func(env.Namespace, ...data.Value) data.Value

var macroType = types.MakeBasic("macro")

// Type makes Call a typed value
func (Call) Type() types.Type {
	return macroType
}

// Equal compares this Call to another for equality
func (Call) Equal(data.Value) bool {
	return false
}

func (c Call) String() string {
	return data.DumpString(c)
}
