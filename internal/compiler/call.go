package compiler

import (
	"github.com/kode4food/ale"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/compiler/encoder"
	"github.com/kode4food/ale/internal/types"
)

// Call represents a code-generating function for the compiler
type Call func(encoder.Encoder, ...ale.Value) error

var (
	CallType = types.MakeBasic("special")

	// compile-time checks for interface implementation
	_ interface {
		data.Mapped
		ale.Typed
	} = Call(nil)
)

// Type makes Call a typed value
func (c Call) Type() ale.Type {
	return types.MakeLiteral(CallType, c)
}

// Equal makes Call a typed Value
func (Call) Equal(ale.Value) bool {
	return false
}

func (c Call) Get(key ale.Value) (ale.Value, bool) {
	return data.DumpMapped(c).Get(key)
}
