package basic

import "github.com/kode4food/ale/types"

type (
	// AnyType accepts a Value of any other Type
	AnyType interface {
		types.Basic
		any() // marker
	}

	any struct {
		types.Basic
	}
)

// Any accepts a Value of any other Type. This will be the default Type
// annotation for data structures
var Any AnyType = &any{
	Basic: New("any"),
}

func (*any) any() {}

func (*any) Accepts(types.Checker, types.Type) bool {
	return true
}
