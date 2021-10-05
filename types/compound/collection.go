package compound

import (
	"fmt"

	"github.com/kode4food/ale/types"
	"github.com/kode4food/ale/types/basic"
	"github.com/kode4food/ale/types/extended"
)

type (
	// CollectionType describes a typed Sequence
	CollectionType interface {
		types.Extended
		collection() // marker
		Element() types.Type
	}

	collection struct {
		types.Extended
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

func makeCollection(base types.Type, elem types.Type) CollectionType {
	return &collection{
		Extended: extended.New(base),
		elem:     elem,
	}
}

func (*collection) collection() {}

func (c *collection) Element() types.Type {
	return c.elem
}

func (c *collection) Name() string {
	return fmt.Sprintf("%s of %s", c.Extended.Name(), c.elem.Name())
}

func (c *collection) Accepts(other types.Type) bool {
	if c == other {
		return true
	}
	if other, ok := other.(CollectionType); ok {
		return c.Extended.Accepts(other) &&
			c.elem.Accepts(other.Element())
	}
	return false
}
