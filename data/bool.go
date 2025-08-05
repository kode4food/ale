package data

import (
	"math/rand/v2"

	"github.com/kode4food/ale"
	"github.com/kode4food/ale/internal/lang"
	"github.com/kode4food/ale/internal/types"
)

// Bool represents the data True or False
type Bool bool

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
func (b Bool) Type() ale.Type {
	return types.MakeLiteral(types.BasicBoolean, b)
}

// HashCode returns the hash code for this Bool
func (b Bool) HashCode() uint64 {
	if b {
		return trueHash
	}
	return falseHash
}
