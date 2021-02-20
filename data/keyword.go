package data

import "math/rand"

// Keyword is a Value that represents a Name that resolves to itself
type Keyword Name

var keywordHash = rand.Uint64()

// Name returns the name of the Keyword
func (k Keyword) Name() Name {
	return Name(k)
}

// Call turns Keyword into a Function
func (k Keyword) Call(args ...Value) Value {
	m := args[0].(Mapped)
	res, ok := m.Get(k)
	if !ok && len(args) > 1 {
		return args[1]
	}
	return res
}

// Convention returns the Function's calling convention
func (k Keyword) Convention() Convention {
	return ApplicativeCall
}

// CheckArity performs a compile-time arity check for the Function
func (k Keyword) CheckArity(argCount int) error {
	return checkRangedArity(1, 2, argCount)
}

// Equal compares this Keyword to another for equality
func (k Keyword) Equal(v Value) bool {
	if v, ok := v.(Keyword); ok {
		return k == v
	}
	return false
}

// String converts Keyword into a string
func (k Keyword) String() string {
	return ":" + string(k)
}

// HashCode returns the hash code for this Keyword
func (k Keyword) HashCode() uint64 {
	return keywordHash * HashString(string(k))
}
