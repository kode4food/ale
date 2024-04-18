package macro

import (
	"github.com/kode4food/ale/internal/types"
	"github.com/kode4food/ale/pkg/data"
	"github.com/kode4food/ale/pkg/env"
)

// Call represents a macro's calling signature
type Call func(env.Namespace, ...data.Value) data.Value

var CallType = types.MakeBasic("macro")

// Type makes Call a typed value
func (Call) Type() types.Type {
	return CallType
}

// Equal compares this Call to another for equality
func (Call) Equal(data.Value) bool {
	return false
}

func (c Call) Get(key data.Value) (data.Value, bool) {
	return data.DumpMapped(c).Get(key)
}
