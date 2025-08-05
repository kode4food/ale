package types

import (
	"fmt"

	"github.com/kode4food/ale"
)

// Object describes a typed set of Key/Value Pairs
type Object struct {
	Basic
	key   ale.Type
	value ale.Type
}

// MakeObject declares a new ObjectType that will only allow keys and values of
// the provided types
func MakeObject(key ale.Type, value ale.Type) ale.Type {
	return &Object{
		Basic: BasicObject,
		key:   key,
		value: value,
	}
}

func (o *Object) Key() ale.Type {
	return o.key
}

func (o *Object) Value() ale.Type {
	return o.value
}

func (o *Object) Name() string {
	return fmt.Sprintf("%s(%s->%s)",
		o.Basic.Name(), o.key.Name(), o.value.Name(),
	)
}

func (o *Object) Accepts(other ale.Type) bool {
	if other, ok := other.(*Object); ok {
		return o == other || compoundAccepts(o, other)
	}
	return false
}

func (o *Object) accepts(c *checker, other ale.Type) bool {
	if other, ok := other.(*Object); ok {
		return o == other ||
			o.Basic.Accepts(other.Basic) &&
				c.acceptsChild(o.key, other.Key()) &&
				c.acceptsChild(o.value, other.Value())
	}
	return false
}

func (o *Object) Equal(other ale.Type) bool {
	if other, ok := other.(*Object); ok {
		return o == other ||
			o.Basic.Equal(other.Basic) &&
				o.key.Equal(other.key) &&
				o.value.Equal(other.value)
	}
	return false
}
