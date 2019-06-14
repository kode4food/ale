package data

import (
	"bytes"
	"fmt"
)

// Object maps values to values and supports the universal design pattern
type Object map[Value]Value

// Error messages
const (
	ObjectNotPaired = "object does not contain an even number of elements"
	ValueNotFound   = "value not found in object: %s"
)

// Standard Keys
const (
	TypeKey     = Keyword("type")
	InstanceKey = Keyword("instance")
)

// NewObject instantiates a new Object instance
func NewObject(args ...Value) Object {
	l := len(args)
	if l%2 != 0 {
		panic(fmt.Errorf(ObjectNotPaired))
	}
	res := Object{}
	for i := len(args) - 2; i >= 0; i -= 2 {
		res[args[i]] = args[i+1]
	}
	return res
}

// Get attempts to retrieve a Value from an Object
func (o Object) Get(k Value) (Value, bool) {
	if v, ok := o[k]; ok {
		return v, ok
	}
	return Null, false
}

// Merge creates a new Object that is the result of merging this and another
func (o Object) Merge(v Object) Object {
	res := o.Copy()
	for k, v := range v {
		res[k] = v
	}
	return res
}

// MustGet retrieves a Value from an Object or explodes
func (o Object) MustGet(k Value) Value {
	if v, ok := o.Get(k); ok {
		return v
	}
	panic(fmt.Errorf(ValueNotFound, k))
}

// Copy creates an exact copy of the current Object
func (o Object) Copy() Object {
	newProps := make(Object, len(o))
	for k, v := range o {
		newProps[k] = v
	}
	return newProps
}

// Caller turns Object into a callable type
func (o Object) Caller() Call {
	return makeMappedCall(o)
}

// Convention returns the function's calling convention
func (o Object) Convention() Convention {
	return ApplicativeCall
}

// CheckArity performs a compile-time arity check for the function
func (o Object) CheckArity(argCount int) error {
	return checkRangedArity(1, 2, argCount)
}

// First returns the first pair of this Object
func (o Object) First() Value {
	return o.toSeq().First()
}

// Rest returns the rest of the pairs of this Object
func (o Object) Rest() Sequence {
	return o.toSeq().Rest()
}

// Split performs a sequencing split of the pairs of this Object
func (o Object) Split() (Value, Sequence, bool) {
	return o.toSeq().Split()
}

// IsEmpty returns whether or not this Object has no pairs
func (o Object) IsEmpty() bool {
	return len(o) == 0
}

// Count returns the number of pairs in this Object
func (o Object) Count() int {
	return len(o)
}

func (o Object) toSeq() List {
	var res List = EmptyList
	for k, v := range o {
		res = res.Prepend(NewCons(k, v)).(List)
	}
	return res
}

// String converts this Value into a string
func (o Object) String() string {
	var buf bytes.Buffer
	buf.WriteString("{")
	idx := 0
	for k, v := range o {
		if idx > 0 {
			buf.WriteString(" ")
		}
		buf.WriteString(MaybeQuoteString(k))
		buf.WriteString(" ")
		buf.WriteString(MaybeQuoteString(v))
		idx++
	}
	buf.WriteString("}")
	return buf.String()
}
