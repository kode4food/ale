package types

import "github.com/kode4food/ale/data"

type (
	// Type is the basic interface for type checking
	Type interface {
		Name() data.Name
		Satisfies(Type) error
	}

	// Boolean represents the boolean native type
	Boolean interface {
		Type
		Boolean()
	}

	// Numeric represents the numeric native type
	Numeric interface {
		Type
		Numeric()
	}

	// String represents the string native type
	String interface {
		Type
		String()
	}

	// Nil represents the empty list or nil value as a distinct type
	Nil interface {
		Type
		Nil()
	}

	// Composite represents a composite type (list, tuple, record, sum)
	Composite interface {
		Type
		Composite()
	}

	// List represents a list type, where each element is the same type
	List interface {
		Composite
		List()
		Type() Type
	}

	// Tuple represents a fixed set of independently typed values
	Tuple interface {
		Composite
		Tuple()
		Types() []Type
	}

	// Record represents a data structure with named fields
	Record interface {
		Composite
		Record()
		Entries() []struct {
			Name data.Name
			Type Type
		}
	}

	// Sum represents a union of types
	Sum interface {
		Composite
		Sum()
		Types() []Type
	}
)
