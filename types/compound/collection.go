package compound

import (
	"fmt"

	"github.com/kode4food/ale/types"
	"github.com/kode4food/ale/types/basic"
)

type (
	// CollectionType describes a typed Sequence
	CollectionType interface {
		types.BasicType
		collection() // marker
		Element() types.Type
	}

	collection struct {
		types.BasicType
		elem types.Type
	}
)

// List declares a new CollectionType that will only allow a List with
// elements of the provided elem Type
func List(elem types.Type) CollectionType {
	return makeCollection(basic.List, elem)
}

// Vector declares a new CollectionType that will only allow a Vector with
// elements of the provided elem Type
func Vector(elem types.Type) CollectionType {
	return makeCollection(basic.Vector, elem)
}

func makeCollection(base types.BasicType, elem types.Type) CollectionType {
	return &collection{
		BasicType: base,
		elem:      elem,
	}
}

func (*collection) collection() {}

func (c *collection) Element() types.Type {
	return c.elem
}

func (c *collection) Name() string {
	return fmt.Sprintf("%s of %s", c.BasicType.Name(), c.elem.Name())
}

func (c *collection) Accepts(other types.Type) bool {
	if c == other {
		return true
	}
	if other, ok := other.(CollectionType); ok {
		return c.BasicType.Accepts(other) &&
			c.elem.Accepts(other.Element())
	}
	return false
}
