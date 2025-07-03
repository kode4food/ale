package data

import (
	"math/rand/v2"

	"github.com/kode4food/ale/internal/types"
)

type (
	// Bool represents the data True or False
	Bool bool

	// Name represents a Value's name. Not itself a Value
	Name string

	// Value is the generic interface for all data
	Value interface {
		Equal(Value) bool
	}

	// Named is the generic interface for values that are named
	Named interface {
		Value
		Name() Name
	}

	// Typed is the generic interface for values that are typed
	Typed interface {
		Value
		Type() types.Type
	}

	// Mapped is the interface for Values that have accessible properties
	Mapped interface {
		Value
		Get(Value) (Value, bool)
	}
)

const (
	// True represents the boolean value of True
	True Bool = true

	// TrueLiteral represents the literal value of True
	TrueLiteral = "#t"

	// False represents the boolean value of false
	False Bool = false

	// FalseLiteral represents the literal value of False
	FalseLiteral = "#f"
)

var (
	trueHash  = rand.Uint64()
	falseHash = rand.Uint64()

	// compile-time checks for interface implementation
	_ interface {
		Hashed
		Procedure
		Typed
	} = False
)

func Equal(l, r Value) bool {
	return l.Equal(r)
}

// Call returns the Bool value regardless of the arguments passed
func (b Bool) Call(...Value) Value {
	return b
}

func (b Bool) CheckArity(int) error {
	return nil
}

// Equal compares this Bool to another for equality
func (b Bool) Equal(other Value) bool {
	return b == other
}

// String converts this Value into a string
func (b Bool) String() string {
	if b {
		return TrueLiteral
	}
	return FalseLiteral
}

// Type returns the Type for this Bool Value
func (Bool) Type() types.Type {
	return types.BasicBoolean
}

// HashCode returns the hash code for this Bool
func (b Bool) HashCode() uint64 {
	if b {
		return trueHash
	}
	return falseHash
}
