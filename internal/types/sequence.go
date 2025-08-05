package types

import (
	"fmt"

	"github.com/kode4food/ale"
)

// MakeListOf declares a new PairType that will only allow a BasicList with
// elements of the provided elem Type
func MakeListOf(elem ale.Type) ale.Type {
	return sequenceOf(BasicList, elem)
}

// MakeVectorOf declares a new VectorType that will only allow a BasicVector
// with elements of the provided elem Type
func MakeVectorOf(elem ale.Type) ale.Type {
	return sequenceOf(BasicVector, elem)
}

func sequenceOf(base *basic, elem ale.Type) ale.Type {
	name := fmt.Sprintf("%s(%s)", base.Name(), elem.Name())
	first := &Pair{
		Basic: base,
		name:  name,
		car:   elem,
	}
	rest := &Union{
		name:    name,
		Basic:   base,
		options: typeList{first, BasicNull},
	}
	first.cdr = rest
	return rest
}
