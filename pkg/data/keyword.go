package data

import (
	"math/rand/v2"

	"github.com/kode4food/ale/internal/lang"
	"github.com/kode4food/ale/internal/types"
)

// Keyword is a Value that represents a Name that resolves to itself
type Keyword string

var kwdSalt = rand.Uint64()

// Name returns the name of the Keyword
func (k Keyword) Name() Local {
	return Local(k)
}

// Call turns Keyword into a Caller
func (k Keyword) Call(args ...Value) Value {
	m := args[0].(Mapped)
	res, ok := m.Get(k)
	if !ok && len(args) > 1 {
		return args[1]
	}
	return res
}

// CheckArity performs a compile-time arity check for the Keyword
func (k Keyword) CheckArity(argc int) error {
	return CheckRangedArity(1, 2, argc)
}

// Equal compares this Keyword to another for equality
func (k Keyword) Equal(other Value) bool {
	return k == other
}

// String converts Keyword into a string
func (k Keyword) String() string {
	return lang.KwdPrefix + string(k)
}

// Type returns the Type for this Keyword Value
func (Keyword) Type() types.Type {
	return types.BasicKeyword
}

// HashCode returns the hash code for this Keyword
func (k Keyword) HashCode() uint64 {
	return kwdSalt ^ HashString(string(k))
}
