package macro

import (
	"github.com/kode4food/ale"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/env"
	"github.com/kode4food/ale/internal/types"
)

// Call represents a macro's calling signature
type Call func(env.Namespace, ...ale.Value) ale.Value

var (
	CallType = types.MakeBasic("macro")

	// compile-time checks for interface implementation
	_ interface {
		data.Mapped
		ale.Typed
	} = Call(nil)
)

// Type makes Call a typed value
func (Call) Type() ale.Type {
	return CallType
}

// Equal compares this Call to another for equality
func (Call) Equal(ale.Value) bool {
	return false
}

func (c Call) Get(key ale.Value) (ale.Value, bool) {
	return data.DumpMapped(c).Get(key)
}
