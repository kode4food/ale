package types

import (
	"fmt"

	"github.com/kode4food/ale"
)

// Object describes a typed set of Key/Value Pairs
type Object struct {
	basic
	key   ale.Type
	value ale.Type
}

// MakeObject declares a new ObjectType that will only allow keys and values of
// the provided types
func MakeObject(key ale.Type, value ale.Type) ale.Type {
	return &Object{
		basic: BasicObject,
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
		o.basic.Name(), o.key.Name(), o.value.Name(),
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
			o.basic.Accepts(other.basic) &&
				c.acceptsChild(o.key, other.Key()) &&
				c.acceptsChild(o.value, other.Value())
	}
	return false
}

func (o *Object) Equal(other ale.Type) bool {
	if other, ok := other.(*Object); ok {
		return o == other ||
			o.basic.Equal(other.basic) &&
				o.key.Equal(other.key) &&
				o.value.Equal(other.value)
	}
	return false
}
