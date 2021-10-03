package types

// Type describes the type compatibility for a Value
type Type interface {
	// Name identifies this Type
	Name() string

	// Accepts determined if this Type will accept the provided Type for
	// binding. This will generally mean that the provided Type satisfies
	// the contract of the receiver.
	Accepts(Type) bool
}
