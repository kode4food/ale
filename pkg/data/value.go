package data

import (
	"math/rand/v2"

	"github.com/kode4food/ale/internal/types"
)

type (
	// Bool represents the data True or False
	Bool bool

	// Value is the generic interface for all data
	Value interface {
		Equal(Value) bool
	}

	// Typed is the generic interface for data that are typed
	Typed interface {
		Type() types.Type
	}

	// Counted interfaces allow a Value to return a count of its items
	Counted interface {
		Count() int
	}

	// Indexed is the interface for data that have indexed elements
	Indexed interface {
		ElementAt(int) (Value, bool)
	}

	// Mapped is the interface for Values that have accessible properties
	Mapped interface {
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
)

func Equal(l, r Value) bool {
	return l.Equal(r)
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
