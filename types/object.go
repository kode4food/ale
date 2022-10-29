package types

import "fmt"

type (
	// ObjectType describes a typed set of Key/Value Pairs
	ObjectType interface {
		Type
		object() // marker
		Key() Type
		Value() Type
		Accepts(*Checker, Type) bool
	}

	object struct {
		BasicType
		key   Type
		value Type
	}
)

// Object declares a new ObjectType that will only allow keys and values of the
// provided types
func Object(key Type, value Type) ObjectType {
	return &object{
		BasicType: AnyObject,
		key:       key,
		value:     value,
	}
}

func (*object) object() {}

func (o *object) Key() Type {
	return o.key
}

func (o *object) Value() Type {
	return o.value
}

func (o *object) Name() string {
	return fmt.Sprintf("%s(%s->%s)",
		o.BasicType.Name(), o.key.Name(), o.value.Name(),
	)
}

func (o *object) Accepts(c *Checker, other Type) bool {
	if other, ok := other.(ObjectType); ok {
		return other.IsA(o.BasicType) &&
			c.AcceptsChild(o.key, other.Key()) &&
			c.AcceptsChild(o.value, other.Value())
	}
	return false
}
