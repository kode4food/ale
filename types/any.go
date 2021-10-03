package types

type (
	// AnyType accepts a Value of any other Type
	AnyType interface {
		Type
		any() // marker
	}

	any struct{}
)

// Any accepts a Value of any other Type. This will be the default Type
// annotation for data structures
var Any AnyType = &any{}

func (*any) any() {}

func (*any) Name() string {
	return "any"
}

func (*any) Accepts(Type) bool {
	return true
}
