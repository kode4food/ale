package compound

import "github.com/kode4food/ale/types"

type (
	// TupleType describes a fixed-length List or Vector, with a specified
	// set of positional element Types
	TupleType interface {
		types.Type
		tuple() // marker
		Elements() []types.Type
	}

	tuple struct {
		elems []types.Type
	}
)

// Tuple declares a new TupleType that will only allow a List or Vector
// with positional elements of the provided Types
func Tuple(elems ...types.Type) TupleType {
	return &tuple{
		elems: elems,
	}
}

func (*tuple) tuple() {}

func (*tuple) Name() string {
	return "tuple"
}

func (t *tuple) Elements() []types.Type {
	return t.elems
}

func (t *tuple) Accepts(other types.Type) bool {
	if t == other {
		return true
	}
	if other, ok := other.(TupleType); ok {
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
