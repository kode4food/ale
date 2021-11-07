package compound

import (
	"fmt"

	"github.com/kode4food/ale/types"
	"github.com/kode4food/ale/types/basic"
	"github.com/kode4food/ale/types/extended"
)

type (
	// VectorType describes a typed fixed array of Values
	VectorType interface {
		types.Extended
		collection() // marker
		Element() types.Type
	}

	vector struct {
		types.Extended
		elem types.Type
	}
)

// Vector declares a new VectorType that will only allow a Vector with elements
// of the provided elem Type
func Vector(elem types.Type) VectorType {
	return &vector{
		Extended: extended.New(basic.Vector),
		elem:     elem,
	}
}

func (*vector) collection() {}

func (v *vector) Element() types.Type {
	return v.elem
}

func (v *vector) Name() string {
	return fmt.Sprintf("%s(%s)", v.Extended.Name(), v.elem.Name())
}

func (v *vector) Accepts(c types.Checker, other types.Type) bool {
	if v == other {
		return true
	}
	if other, ok := other.(VectorType); ok {
		return c.Check(v.Extended).Accepts(other) != nil &&
			c.Check(v.elem).Accepts(other.Element()) != nil
	}
	return false
}
