package data

import "fmt"

type (
	// Any makes "interface{}" less ugly
	Any = interface{}

	// Bool represents the values True or False
	Bool bool

	// Value is the generic interface for all values
	Value interface {
		fmt.Stringer
	}

	// Values represents a set of Values
	Values []Value

	// Name is a Variable name
	Name string

	// Names represents a set of Names
	Names []Name

	// Typed is the generic interface for values that are typed
	Typed interface {
		Type() Name
	}

	// Named is the generic interface
	Named interface {
		Name() Name
	}

	// Counted interfaces allow a Value to return a count of its items
	Counted interface {
		Count() int
	}

	// Mapped is the interface for values that have retrievable properties
	Mapped interface {
		Get(Value) (Value, bool)
	}

	// Indexed is the interface for values that have indexed elements
	Indexed interface {
		ElementAt(int) (Value, bool)
	}
)

const (
	// True represents the boolean value of True
	True Bool = true

	// False represents the boolean value of false
	False Bool = false
)

func (v Values) String() string {
	return DumpString(v)
}

// Name makes Name Named
func (n Name) Name() Name {
	return n
}

// String converts this Value into a string
func (n Name) String() string {
	return string(n)
}

// String converts this Value into a string
func (b Bool) String() string {
	if bool(b) {
		return "true"
	}
	return "false"
}

// Truthy evaluates whether or not a Value is truthy
func Truthy(v Value) bool {
	if v == False || v == Null {
		return false
	}
	return true
}
