package compound

import (
	"fmt"

	"github.com/kode4food/ale/types"
	"github.com/kode4food/ale/types/basic"
	"github.com/kode4food/ale/types/extended"
)

type (
	// TupleType describes a fixed-length List or Vector, with a specified
	// set of positional element Types
	TupleType interface {
		types.Extended
		tuple() // marker
		Elements() []types.Type
	}

	tuple struct {
		types.Extended
		elems typeList
	}
)

var tupleAcceptor = Union(basic.List, basic.Vector)

// Tuple declares a new TupleType that will only allow a List or Vector
// with positional elements of the provided Types
func Tuple(elems ...types.Type) TupleType {
	return makeTuple(tupleAcceptor, elems)
}

// ListTuple declares a new TupleType that only accepts a List as its base
func ListTuple(elems ...types.Type) TupleType {
	return makeTuple(basic.List, elems)
}

// VectorTuple declares a new TupleType that only accepts a Vector as its base
func VectorTuple(elems ...types.Type) TupleType {
	return makeTuple(basic.Vector, elems)
}

func makeTuple(base types.Type, elems []types.Type) TupleType {
	return &tuple{
		Extended: extended.New(base),
		elems:    elems,
	}
}

func (*tuple) tuple() {}

func (t *tuple) Name() string {
	return fmt.Sprintf("tuple(%s)", t.elems.name())
}

func (t *tuple) Elements() []types.Type {
	return t.elems
}

func (t *tuple) Accepts(other types.Type) bool {
	if t == other {
		return true
	}
	if other, ok := other.(TupleType); ok {
		if !t.Extended.Accepts(other.Base()) {
			return false
		}
		oe := other.Elements()
		if len(t.elems) != len(oe) {
			return false
		}
		for i, elem := range t.elems {
			if !elem.Accepts(oe[i]) {
				return false
			}
		}
		return true
	}
	return false
}
