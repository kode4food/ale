package data

import (
	"bytes"
	"fmt"
)

// Object maps values to values and supports the universal design pattern
type Object map[Value]Value

// Error messages
const (
	ValueNotFound = "value not found in object: %s"
)

const prototypeKey = Keyword("prototype")

// GetPrototype returns the Object's prototype (lookup chain)
func (o Object) GetPrototype() (Object, bool) {
	if v, ok := o[prototypeKey]; ok {
		if proto, ok := v.(Object); ok {
			return proto, true
		}
	}
	return nil, false
}

// Get attempts to retrieve a Value from an Object
func (o Object) Get(k Value) (Value, bool) {
	if v, ok := o[k]; ok {
		return v, ok
	}
	if proto, ok := o.GetPrototype(); ok {
		return proto.Get(k)
	}
	return Nil, false
}

// MustGet retrieves a Value from an Object or explodes
func (o Object) MustGet(k Value) Value {
	if v, ok := o.Get(k); ok {
		return v
	}
	panic(fmt.Errorf(ValueNotFound, k))
}

// Extend instantiates a new Object using the current Object as a prototype
func (o Object) Extend(properties Object) Object {
	newObject := properties.Flatten().Copy()
	newObject[prototypeKey] = o
	return newObject
}

// Copy creates an exact copy of the current Object
func (o Object) Copy() Object {
	newProps := make(Object, len(o))
	for k, v := range o {
		newProps[k] = v
	}
	return newProps
}

// Flatten returns completely flattened copy of an Object, removing the entire
// prototype chain from the result
func (o Object) Flatten() Object {
	if proto, ok := o.GetPrototype(); ok {
		pf := proto.Flatten()
		r := make(Object, len(pf)+len(o))
		for k, v := range pf {
			r[k] = v
		}
		for k, v := range o {
			if k == prototypeKey {
				continue
			}
			r[k] = v
		}
		return r
	}
	return o
}

// Caller turns Object into a callable type
func (o Object) Caller() Call {
	return makeMappedCall(o)
}

// String converts this Value into a string
func (o Object) String() string {
	var buf bytes.Buffer
	buf.WriteString("{")
	idx := 0
	for k, v := range o.Flatten() {
		if idx > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(MaybeQuoteString(k))
		buf.WriteString(" ")
		buf.WriteString(MaybeQuoteString(v))
		idx++
	}
	buf.WriteString("}")
	return buf.String()
}
