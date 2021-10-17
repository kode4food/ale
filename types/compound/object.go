package compound

import (
	"fmt"

	"github.com/kode4food/ale/types"
	"github.com/kode4food/ale/types/basic"
	"github.com/kode4food/ale/types/extended"
)

type (
	// ObjectType describes a typed set of Key/Value Pairs
	ObjectType interface {
		types.Extended
		object() // marker
		Key() types.Type
		Value() types.Type
	}

	object struct {
		types.Extended
		key   types.Type
		value types.Type
	}
)

// Object declares a new ObjectType that will only allow keys and values
// of the provided types
func Object(key types.Type, value types.Type) ObjectType {
	return &object{
		Extended: extended.New(basic.Object),
		key:      key,
		value:    value,
	}
}

func (*object) object() {}

func (o *object) Key() types.Type {
	return o.key
}

func (o *object) Value() types.Type {
	return o.value
}

func (o *object) Name() string {
	return fmt.Sprintf("%s(%s->%s)",
		o.Extended.Name(), o.key.Name(), o.value.Name(),
	)
}

func (o *object) Accepts(c types.Checker, other types.Type) bool {
	if o == other {
		return true
	}
	if other, ok := other.(ObjectType); ok {
		return c.Check(o.Extended).Accepts(other) != nil &&
			c.Check(o.key).Accepts(other.Key()) != nil &&
			c.Check(o.value).Accepts(other.Value()) != nil
	}
	return false
}
