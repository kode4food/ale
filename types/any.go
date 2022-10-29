package types

// AnyType accepts a Value of any other Type
type (
	AnyType interface {
		BasicType
		any() // marker
	}

	anyType struct{ BasicType }
)

var Any = &anyType{
	BasicType: Basic("any"),
}

func (*anyType) any() {}

func (a *anyType) IsA(k BasicType) bool {
	return a.Kind() == k.Kind()
}

func (*anyType) Accepts(*Checker, Type) bool {
	return true
}
