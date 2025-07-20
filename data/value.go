package data

import (
	"math/rand/v2"

	"github.com/kode4food/ale"
	"github.com/kode4food/ale/internal/lang"
	"github.com/kode4food/ale/internal/types"
)

type (
	// Bool represents the data True or False
	Bool bool

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

const (
	// True represents the boolean value of True
	True Bool = true

	// False represents the boolean value of false
	False Bool = false
)

var (
	trueHash  = rand.Uint64()
	falseHash = rand.Uint64()

	// compile-time checks for interface implementation
	_ interface {
		Hashed
		Procedure
		ale.Typed
	} = False
)

func Equal(l, r ale.Value) bool {
	return l.Equal(r)
}

// Call returns the Bool value regardless of the arguments passed
func (b Bool) Call(...ale.Value) ale.Value {
	return b
}

func (b Bool) CheckArity(int) error {
	return nil
}

// Equal compares this Bool to another for equality
func (b Bool) Equal(other ale.Value) bool {
	return b == other
}

// String converts this Value into a string
func (b Bool) String() string {
	if b {
		return lang.TrueLiteral
	}
	return lang.FalseLiteral
}

// Type returns the Type for this Bool Value
func (Bool) Type() ale.Type {
	return types.BasicBoolean
}

// HashCode returns the hash code for this Bool
func (b Bool) HashCode() uint64 {
	if b {
		return trueHash
	}
	return falseHash
}
