package data

import "github.com/kode4food/ale"

type (
	// Name represents a Value's name. Not itself a Value
	Name string

	// Named is the generic interface for values that are named
	Named interface {
		ale.Value
		Name() Name
	}

	// Typed is the generic interface for values that are typed

	// Mapped is the interface for Values that have accessible properties
	Mapped interface {
		ale.Value
		Get(ale.Value) (ale.Value, bool)
	}
)

func Equal(l, r ale.Value) bool {
	return l.Equal(r)
}
