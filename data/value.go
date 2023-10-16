package data

import (
	"fmt"
	"hash/maphash"
	"math/rand"

	"github.com/kode4food/ale/internal/types"
)

type (
	// Bool represents the values True or False
	Bool bool

	// Value is the generic interface for all values
	Value interface {
		fmt.Stringer
		Equal(Value) bool
	}

	// Values represent a set of Values
	Values []Value

	// Typed is the generic interface for values that are typed
	Typed interface {
		Type() types.Type
	}

	// Counted interfaces allow a Value to return a count of its items
	Counted interface {
		Count() int
	}

	// Indexed is the interface for values that have indexed elements
	Indexed interface {
		ElementAt(int) (Value, bool)
	}

	// Mapped is the interface for Values that have accessible properties
	Mapped interface {
		Get(Value) (Value, bool)
	}

	// Valuer can return its data as a slice of Values
	Valuer interface {
		Values() Values
	}

	// Hashed can return a hash code for the value
	Hashed interface {
		HashCode() uint64
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
	seed = maphash.MakeSeed()

	trueHash  = rand.Uint64()
	falseHash = rand.Uint64()
)

// Equal compares this Bool to another for equality
func (b Bool) Equal(v Value) bool {
	if v, ok := v.(Bool); ok {
		return b == v
	}
	return false
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

// HashCode returns a hash code for the provided Value. If the Value implements
// the Hashed interface, it will call us the HashCode() method. Otherwise, it
// will create a hash code from the stringified form of the Value
func HashCode(v Value) uint64 {
	if h, ok := v.(Hashed); ok {
		return h.HashCode()
	}
	return HashString(v.String())
}

// HashString returns a hash code for the provided string
func HashString(s string) uint64 {
	var b maphash.Hash
	b.SetSeed(seed)
	_, _ = b.WriteString(s)
	return b.Sum64()
}
