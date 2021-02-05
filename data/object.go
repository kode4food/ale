package data

import (
	"bytes"
	"errors"
	"fmt"
	"sort"
)

// Object maps values to values and supports the universal design pattern
type Object map[Value]Value

// Standard Keys
const (
	NameKey     = Keyword("name")
	TypeKey     = Keyword("type")
	CountKey    = Keyword("count")
	InstanceKey = Keyword("instance")
)

// Error messages
const (
	ErrMapNotPaired  = "map does not contain an even number of elements"
	ErrValueNotFound = "value not found in object: %s"
)

// EmptyObject represents an empty Object
var EmptyObject = Object{}

// NewObject instantiates a new Object instance
func NewObject(pairs ...Pair) Object {
	res := Object{}
	for _, p := range pairs {
		res[p.Car()] = p.Cdr()
	}
	return res
}

// ValuesToObject interprets a set of Values as an Object
func ValuesToObject(v ...Value) (Object, error) {
	if len(v)%2 != 0 {
		return nil, errors.New(ErrMapNotPaired)
	}
	var p Pairs
	for i := len(v) - 2; i >= 0; i -= 2 {
		p = append(p, NewCons(v[i], v[i+1]))
	}
	return NewObject(p...), nil
}

// Get attempts to retrieve a Value from an Object
func (o Object) Get(k Value) (Value, bool) {
	if v, ok := o[k]; ok {
		return v, ok
	}
	return Nil, false
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
	panic(fmt.Errorf(ErrValueNotFound, k))
}

// Copy creates an exact copy of the current Object
func (o Object) Copy() Object {
	newProps := make(Object, len(o))
	for k, v := range o {
		newProps[k] = v
	}
	return newProps
}

// Call turns Object into a Function
func (o Object) Call(args ...Value) Value {
	return mappedCall(o, args)
}

// Convention returns the Function's calling convention
func (o Object) Convention() Convention {
	return ApplicativeCall
}

// CheckArity performs a compile-time arity check for the Function
func (o Object) CheckArity(argCount int) error {
	return checkRangedArity(1, 2, argCount)
}

// First returns the first pair of this Object
func (o Object) First() Value {
	return o.firstFrom(o.sortedKeys())
}

func (o Object) firstFrom(keys []Value) Value {
	if len(keys) > 0 {
		k0 := keys[0]
		return NewCons(k0, o[k0])
	}
	return Nil
}

// Rest returns the rest of the pairs of this Object
func (o Object) Rest() Sequence {
	return o.restFrom(o.sortedKeys())
}

func (o Object) restFrom(keys []Value) Object {
	if len(keys) > 1 {
		rest := make(Object, len(keys)-1)
		for _, k := range keys[1:] {
			rest[k] = o[k]
		}
		return rest
	}
	return EmptyObject
}

// Split performs a sequencing split of the pairs of this Object
func (o Object) Split() (Value, Sequence, bool) {
	if len(o) > 0 {
		keys := o.sortedKeys()
		first := o.firstFrom(keys)
		rest := o.restFrom(keys)
		return first, rest, true
	}
	return Nil, EmptyObject, false
}

func (o Object) sortedKeys() []Value {
	keys := make([]Value, 0, len(o))
	for k := range o {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(l, r int) bool {
		return fmt.Sprintf("%p", keys[l]) < fmt.Sprintf("%p", keys[r])
	})
	return keys
}

// IsEmpty returns whether this Object has no pairs
func (o Object) IsEmpty() bool {
	return len(o) == 0
}

// Count returns the number of pairs in this Object
func (o Object) Count() int {
	return len(o)
}

// Equal compares this Object to another for equality
func (o Object) Equal(v Value) bool {
	if ro, ok := v.(Object); ok {
		if len(o) != len(ro) {
			return false
		}
		for leftKey, leftVal := range o {
			rightVal, ok := ro[leftKey]
			if !ok || !leftVal.Equal(rightVal) {
				return false
			}
		}
		return true
	}
	return false
}

// String converts this Value into a string
func (o Object) String() string {
	var buf bytes.Buffer
	buf.WriteString("{")
	for i, k := range o.sortedKeys() {
		if i > 0 {
			buf.WriteString(" ")
		}
		buf.WriteString(MaybeQuoteString(k))
		buf.WriteString(" ")
		buf.WriteString(MaybeQuoteString(o[k]))
	}
	buf.WriteString("}")
	return buf.String()
}
