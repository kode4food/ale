package types

type (
	// Type describes the type compatibility for a Value
	Type interface {
		// Name identifies this Type
		Name() string

		// Accepts determined if this Type will accept the provided Type for
		// binding. This will generally mean that the provided Type satisfies
		// the contract of the receiver.
		Accepts(Type) bool
	}

	// Basic describes an atomic Type that exposes a comparable Kind
	Basic interface {
		Type

		// Kind returns the Kind for this Type
		Kind() Kind
	}

	// Kind uniquely identifies a Type within a process
	Kind [16]byte

	// Extended describes a Type that extends another Type
	Extended interface {
		Type

		// Base returns the base Type of this extended Type
		Base() Type
	}
)
