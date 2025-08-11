package ale

type (
	// Value is the generic interface for all data
	Value interface {
		// Equal compares this Value to another for equality
		Equal(Value) bool
	}

	// Type describes the type compatibility for a Value
	Type interface {
		// Name identifies this Type
		Name() string

		// Accepts determines if this Type will accept the provided Type for
		// binding. This will generally mean that the provided Type satisfies
		// the contract of the receiver.
		Accepts(Type) bool

		// Equal determines if the provided Type is an equivalent definition
		Equal(Type) bool
	}

	Typed interface {
		Value

		// Type returns the Type for this Value
		Type() Type
	}
)
