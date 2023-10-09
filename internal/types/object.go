package types

import "fmt"

// Object describes a typed set of Key/Value Pairs
type Object struct {
	basic
	key   Type
	value Type
}

// MakeObject declares a new ObjectType that will only allow keys and values of
// the provided types
func MakeObject(key Type, value Type) *Object {
	return &Object{
		basic: BasicObject,
		key:   key,
		value: value,
	}
}

func (o *Object) Key() Type {
	return o.key
}

func (o *Object) Value() Type {
	return o.value
}

func (o *Object) Name() string {
	return fmt.Sprintf("%s(%s->%s)",
		o.basic.Name(), o.key.Name(), o.value.Name(),
	)
}

func (o *Object) Accepts(c *Checker, other Type) bool {
	if o == other {
		return true
	}
	if other, ok := other.(*Object); ok {
		return o.basic.Accepts(c, other) &&
			c.AcceptsChild(o.key, other.Key()) &&
			c.AcceptsChild(o.value, other.Value())
	}
	return false
}

func (o *Object) Equal(other Type) bool {
	if o == other {
		return true
	}
	if other, ok := other.(*Object); ok {
		return o.basic.Equal(other.basic) &&
			o.key.Equal(other.key) &&
			o.value.Equal(other.value)
	}
	return false
}
